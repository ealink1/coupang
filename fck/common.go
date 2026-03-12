package fck

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type returnJSON struct {
	Data struct {
		Content string `json:"content"`
	} `json:"data"`
}

type ConvertResult struct {
	PNGPath string
}

func maybeGunzip(b []byte) ([]byte, error) {
	if len(b) >= 2 && b[0] == 0x1f && b[1] == 0x8b {
		r, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, fmt.Errorf("create gzip reader: %w", err)
		}
		defer r.Close()

		out, err := io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("gunzip: %w", err)
		}
		return out, nil
	}
	return b, nil
}

func pdfBytesToPNG(ctx context.Context, pdfBytes []byte, outPNGPath string, dpi int) error {
	pdfFile, err := os.CreateTemp("", "coupang_label_*.pdf")
	if err != nil {
		return fmt.Errorf("create temp pdf: %w", err)
	}
	pdfPath := pdfFile.Name()
	defer os.Remove(pdfPath)

	if _, err := pdfFile.Write(pdfBytes); err != nil {
		_ = pdfFile.Close()
		return fmt.Errorf("write temp pdf: %w", err)
	}
	if err := pdfFile.Close(); err != nil {
		return fmt.Errorf("close temp pdf: %w", err)
	}

	if err := pdfToPNG(ctx, pdfPath, outPNGPath, dpi); err != nil {
		return err
	}
	return nil
}

func pdfToPNG(ctx context.Context, pdfPath string, outPNGPath string, dpi int) error {
	if _, err := exec.LookPath("pdftoppm"); err == nil {
		return runPdftoppm(ctx, pdfPath, outPNGPath, dpi)
	}

	if _, err := exec.LookPath("sips"); err == nil {
		return runSips(ctx, pdfPath, outPNGPath)
	}

	return fmt.Errorf("pdftoppm not found in PATH (and sips fallback not found)")
}

func runPdftoppm(ctx context.Context, pdfPath string, outPNGPath string, dpi int) error {
	tmpPrefixFile, err := os.CreateTemp("", "coupang_label_ppm_*")
	if err != nil {
		return fmt.Errorf("create temp prefix: %w", err)
	}
	tmpPrefix := tmpPrefixFile.Name()
	_ = tmpPrefixFile.Close()
	_ = os.Remove(tmpPrefix)

	cmd := exec.CommandContext(ctx, "pdftoppm", "-png", "-r", fmt.Sprintf("%d", dpi), "-singlefile", pdfPath, tmpPrefix)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "TZ=UTC")
	cmd.Cancel = func() error { return cmd.Process.Kill() }

	done := make(chan error, 1)
	go func() { done <- cmd.Run() }()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("pdftoppm convert: %w", err)
		}
	case <-time.After(2 * time.Minute):
		_ = cmd.Process.Kill()
		return fmt.Errorf("pdftoppm timeout")
	}

	generated := tmpPrefix + ".png"
	defer os.Remove(generated)
	if err := copyFile(generated, outPNGPath); err != nil {
		return err
	}
	return nil
}

func runSips(ctx context.Context, pdfPath string, outPNGPath string) error {
	tmpOut, err := os.CreateTemp("", "coupang_label_*.png")
	if err != nil {
		return fmt.Errorf("create temp png: %w", err)
	}
	tmpPNG := tmpOut.Name()
	_ = tmpOut.Close()
	_ = os.Remove(tmpPNG)
	defer os.Remove(tmpPNG)

	cmd := exec.CommandContext(ctx, "sips", "-s", "format", "png", pdfPath, "--out", tmpPNG)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sips convert: %w", err)
	}

	if err := copyFile(tmpPNG, outPNGPath); err != nil {
		return err
	}
	return nil
}

func copyFile(from string, to string) error {
	in, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
		return fmt.Errorf("create out dir: %w", err)
	}

	out, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy file: %w", err)
	}
	return nil
}
