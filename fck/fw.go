package fck

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func ConvertFMJSONToFWPNG(ctx context.Context, jsonPath string, outDir string) (*ConvertResult, error) {
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

	pdfBytes, err := renderHTMLTo10x15PDFBytes(ctx, htmlBytes)
	if err != nil {
		return nil, err
	}

	outPNG := filepath.Join(outDir, "fw.png")
	if err := pdfBytesToPNG(ctx, pdfBytes, outPNG, 300); err != nil {
		return nil, err
	}

	return &ConvertResult{PNGPath: outPNG}, nil
}

