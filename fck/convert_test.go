package fck

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
)

func TestConvert711JSONToPNG(t *testing.T) {
	//if os.Getenv("FCK_RUN_CONVERT") != "1" {
	//	t.Skip("set FCK_RUN_CONVERT=1 to run")
	//}

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("runtime caller failed")
	}
	dir := filepath.Dir(file)

	res, err := Convert711JSONToPNG(context.Background(), filepath.Join(dir, "711.json"), dir)
	if err != nil {
		t.Fatalf("convert failed: %v", err)
	}
	if res == nil || res.PNGPath == "" {
		t.Fatalf("unexpected result: %#v", res)
	}
}

func TestConvertFMJSONToFWPNG(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("runtime caller failed")
	}
	dir := filepath.Dir(file)

	res, err := ConvertFMJSONToFWPNG(context.Background(), filepath.Join(dir, "fm.json"), dir)
	if err != nil {
		t.Fatalf("convert failed: %v", err)
	}
	if res == nil || res.PNGPath == "" {
		t.Fatalf("unexpected result: %#v", res)
	}
}
