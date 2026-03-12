package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	maxKB := flag.Int("maxKB", 500, "")
	overwrite := flag.Bool("overwrite", true, "")
	outDir := flag.String("outDir", "", "")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		files = []string{
			"/Users/bre/workspace/jay/coupang/pic/7533d78adb074e6417c37025f12f4c26.jpg",
		}
	}

	maxBytes := int64(*maxKB) * 1024
	if maxBytes <= 0 {
		fmt.Fprintln(os.Stderr, "invalid -maxKB")
		os.Exit(2)
	}

	for _, f := range files {
		if err := compressToLimit(f, maxBytes, *overwrite, *outDir); err != nil {
			fmt.Fprintf(os.Stderr, "compress failed: %s: %v\n", f, err)
			os.Exit(1)
		}
	}
}

func compressToLimit(path string, maxBytes int64, overwrite bool, outDir string) error {
	original, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if int64(len(original)) <= maxBytes {
		fmt.Printf("%s: %d -> %d bytes (skip)\n", filepath.Base(path), len(original), len(original))
		return nil
	}

	img, err := decodeImage(bytes.NewReader(original))
	if err != nil {
		return err
	}

	type attempt struct {
		data    []byte
		quality int
		scale   float64
	}

	var best *attempt
	srcBounds := img.Bounds()
	srcW, srcH := srcBounds.Dx(), srcBounds.Dy()
	if srcW <= 0 || srcH <= 0 {
		return fmt.Errorf("invalid image bounds")
	}

	scales := []float64{1.0, 0.95, 0.9, 0.85, 0.8, 0.75, 0.7, 0.65, 0.6, 0.55, 0.5, 0.45, 0.4, 0.35, 0.3, 0.25, 0.2}
	qualities := []int{88, 84, 80, 76, 72, 68, 64, 60, 56, 52, 48, 44, 40, 36, 32, 28, 24, 20}

	for _, scale := range scales {
		var scaled image.Image = img
		if scale < 0.999 {
			newW := int(float64(srcW)*scale + 0.5)
			newH := int(float64(srcH)*scale + 0.5)
			if newW < 1 {
				newW = 1
			}
			if newH < 1 {
				newH = 1
			}
			scaled = resizeNearest(img, newW, newH)
		}

		for _, q := range qualities {
			b, err := encodeJPEG(scaled, q)
			if err != nil {
				return err
			}
			if best == nil || len(b) < len(best.data) {
				best = &attempt{data: b, quality: q, scale: scale}
			}
			if int64(len(b)) <= maxBytes {
				return writeOutput(path, b, original, overwrite, outDir, q, scale)
			}
		}
	}

	if best == nil {
		return fmt.Errorf("no compression attempt produced output")
	}

	return fmt.Errorf("unable to reach limit: got %d bytes (best q=%d scale=%.2f)", len(best.data), best.quality, best.scale)
}

func decodeImage(r io.Reader) (image.Image, error) {
	img, format, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	_ = format
	return img, nil
}

func encodeJPEG(img image.Image, quality int) ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(512 * 1024)
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeOutput(srcPath string, data []byte, original []byte, overwrite bool, outDir string, q int, scale float64) error {
	dstPath := srcPath
	if !overwrite {
		ext := filepath.Ext(srcPath)
		base := strings.TrimSuffix(filepath.Base(srcPath), ext)
		dstName := fmt.Sprintf("%s_compressed%s", base, ext)
		if outDir != "" {
			dstPath = filepath.Join(outDir, dstName)
		} else {
			dstPath = filepath.Join(filepath.Dir(srcPath), dstName)
		}
	} else if outDir != "" {
		dstPath = filepath.Join(outDir, filepath.Base(srcPath))
	}

	if err := os.WriteFile(dstPath, data, 0o644); err != nil {
		return err
	}
	fmt.Printf("%s: %d -> %d bytes (q=%d scale=%.2f)\n", filepath.Base(srcPath), len(original), len(data), q, scale)
	return nil
}

func resizeNearest(src image.Image, newW, newH int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	sb := src.Bounds()
	sw := sb.Dx()
	sh := sb.Dy()

	for y := 0; y < newH; y++ {
		sy := sb.Min.Y + (y*sh)/newH
		for x := 0; x < newW; x++ {
			sx := sb.Min.X + (x*sw)/newW
			dst.Set(x, y, src.At(sx, sy))
		}
	}

	return dst
}
