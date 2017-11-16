package diff

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"testing"

	_ "image/jpeg"
)

func TestCompareFiles(t *testing.T) {
	diff, percent, err := CompareFiles(
		"testdata/test-only-text.png",
		"testdata/test-text-number.png",
	)
	// writePng(diff, "testdata/diff.png")

	if percent != 1.37 {
		t.Errorf("CompareFiles() percent = %v, want 1.37", percent)
	}
	if diff == nil {
		t.Errorf("diff should be not nil")
	}
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestCompareFilesWillFailIfNoSuchFile(t *testing.T) {
	_, _, err := CompareFiles(
		"testdata/test-only-text.png",
		"no-such-file.png",
	)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	_, _, err = CompareFiles(
		"no-such-file.png",
		"testdata/test-only-text.png",
	)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestCompareFilesWillFailFileNotImage(t *testing.T) {
	_, _, err := CompareFiles(
		"testdata/test-only-text.png",
		"diff_test.go",
	)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestCompareImages(t *testing.T) {
	white := image.NewRGBA(image.Rect(0, 0, 200, 200))
	draw.Draw(white, white.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)

	black := image.NewRGBA(image.Rect(0, 0, 200, 200))
	draw.Draw(black, black.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	draw.Draw(black, image.Rect(50, 50, 100, 100), &image.Uniform{color.Black}, image.Pt(50, 50), draw.Src)

	diff, percent, err := CompareImages(white, black)

	// writePng(white, "testdata/white.png")
	// writePng(black, "testdata/black.png")
	// writePng(diff, "testdata/diff.png")

	if percent != 6.25 || err != nil || diff == nil {
		t.Fail()
	}
}

func TestCompareImagesWithDiffSizes(t *testing.T) {
	white := image.NewRGBA(image.Rect(0, 0, 200, 200))
	draw.Draw(white, white.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)

	black := image.NewRGBA(image.Rect(0, 0, 100, 100))
	draw.Draw(black, image.Rect(50, 50, 100, 100), &image.Uniform{color.Black}, image.Pt(50, 50), draw.Src)

	_, _, err := CompareImages(white, black)
	if err == nil {
		t.Fail()
	}
}

func TestIsEqualColor(t *testing.T) {
	a := color.RGBA{uint8(255), uint8(255), uint8(255), 255}
	b := color.RGBA{uint8(0), uint8(0), uint8(0), 255}

	if isEqualColor(a, b) {
		t.Fail()
	}
}

func writePng(img image.Image, filename string) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return err
	}
	return ioutil.WriteFile(filename, buf.Bytes(), 0644)
}
