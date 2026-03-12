package fck

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Convert711JSONToPNG(ctx context.Context, jsonPath string, outDir string) (*ConvertResult, error) {
	raw, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("read json: %w", err)
	}

	var payload returnJSON
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(payload.Data.Content)
	if err != nil {
		return nil, fmt.Errorf("base64 decode content: %w", err)
	}

	htmlBytes, err := maybeGunzip(decoded)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, fmt.Errorf("create out dir: %w", err)
	}

	outPNG := filepath.Join(outDir, "711.png")

	pdfBytes, err := renderHTMLTo10x15PDFBytes(ctx, htmlBytes)
	if err != nil {
		return nil, err
	}

	if err := pdfBytesToPNG(ctx, pdfBytes, outPNG, 300); err != nil {
		return nil, err
	}

	return &ConvertResult{PNGPath: outPNG}, nil
}

func renderHTMLTo10x15PDFBytes(ctx context.Context, htmlBytes []byte) ([]byte, error) {
	chromePath, err := findChrome()
	if err != nil {
		return nil, err
	}

	modified := inject10x15PageCSS(string(htmlBytes))

	htmlFile, err := os.CreateTemp("", "coupang_711_*.html")
	if err != nil {
		return nil, fmt.Errorf("create temp html: %w", err)
	}
	htmlPath := htmlFile.Name()
	if _, err := io.Copy(htmlFile, bytes.NewReader([]byte(modified))); err != nil {
		_ = htmlFile.Close()
		_ = os.Remove(htmlPath)
		return nil, fmt.Errorf("write temp html: %w", err)
	}
	if err := htmlFile.Close(); err != nil {
		_ = os.Remove(htmlPath)
		return nil, fmt.Errorf("close temp html: %w", err)
	}
	defer os.Remove(htmlPath)

	pdfFile, err := os.CreateTemp("", "coupang_711_*.pdf")
	if err != nil {
		return nil, fmt.Errorf("create temp pdf: %w", err)
	}
	pdfPath := pdfFile.Name()
	_ = pdfFile.Close()
	_ = os.Remove(pdfPath)
	defer os.Remove(pdfPath)

	fileURL := "file://" + htmlPath
	cmd := exec.CommandContext(ctx, chromePath,
		"--headless",
		"--disable-gpu",
		"--no-first-run",
		"--no-default-browser-check",
		"--print-to-pdf="+pdfPath,
		"--print-to-pdf-no-header",
		fileURL,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("chrome print-to-pdf: %w", err)
	}

	pdfBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("read temp pdf: %w", err)
	}
	return pdfBytes, nil
}

func inject10x15PageCSS(html string) string {
	style := `<style>@page{size:100mm 150mm;margin:0;}html{margin:0;padding:0;width:100%;}body{margin:0 auto !important;padding:0;width:80%;transform:scale(1.25);transform-origin:top center;}</style>`
	lower := strings.ToLower(html)
	idx := strings.Index(lower, "<head>")
	if idx >= 0 {
		return html[:idx+len("<head>")] + style + html[idx+len("<head>"):]
	}
	idx = strings.Index(lower, "<head")
	if idx >= 0 {
		end := strings.Index(lower[idx:], ">")
		if end >= 0 {
			end = idx + end + 1
			return html[:end] + style + html[end:]
		}
	}
	return style + html
}

func findChrome() (string, error) {
	candidates := []string{
		"google-chrome",
		"chromium",
		"chromium-browser",
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
	}

	for _, c := range candidates {
		if strings.HasPrefix(c, "/") {
			if st, err := os.Stat(c); err == nil && !st.IsDir() {
				return c, nil
			}
			continue
		}
		if p, err := exec.LookPath(c); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("chrome/chromium not found in PATH (and no app path found)")
}
