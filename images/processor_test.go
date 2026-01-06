package images

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"testing"
)

// createTestImage creates a simple image for testing.
func createTestImage(width, height int, format string) (io.Reader, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Fill with some color
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{uint8(x % 255), uint8(y % 255), 100, 255})
		}
	}

	var buf bytes.Buffer
	var err error
	switch format {
	case "png":
		err = png.Encode(&buf, img)
	case "jpg", "jpeg":
		err = jpeg.Encode(&buf, img, nil)
	}
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

func TestResizeAspectRatio(t *testing.T) {
	// Create a 200x100 image (2:1 ratio)
	input, err := createTestImage(200, 100, "png")
	if err != nil {
		t.Fatalf("failed to create test image: %v", err)
	}

	processor := NewImageProcessor()
	// Resize to fit in 100x100. Should result in 100x50 to keep 2:1 ratio.
	processor = processor.AddFilter(ResizeAspectRatio(100, 100, true))

	output, err := processor.Process(input)
	if err != nil {
		t.Fatalf("process failed: %v", err)
	}

	// Decode output to check dimensions
	outImg, _, err := image.Decode(output)
	if err != nil {
		t.Fatalf("failed to decode output: %v", err)
	}

	bounds := outImg.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 50 {
		t.Errorf("expected 100x50, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestFormatWebp(t *testing.T) {
	input, err := createTestImage(50, 50, "png")
	if err != nil {
		t.Fatalf("failed to create test image: %v", err)
	}

	processor := NewImageProcessor()
	processor = processor.AddFilter(FormatWebp())

	output, err := processor.Process(input)
	if err != nil {
		t.Fatalf("process failed: %v", err)
	}

	// Let's check header manually for RIFF and WEBP to be sure
	buf, _ := io.ReadAll(output)
	if len(buf) < 12 {
		t.Fatalf("output too short to be webp")
	}
	if string(buf[0:4]) != "RIFF" || string(buf[8:12]) != "WEBP" {
		t.Errorf("header does not look like WEBP: %v", string(buf[0:12]))
	}
}

func TestResizeMaintainFormat(t *testing.T) {
	input, err := createTestImage(50, 50, "jpg")
	if err != nil {
		t.Fatalf("failed to create test image: %v", err)
	}

	processor := NewImageProcessor()
	processor = processor.AddFilter(Resize(25, 25))

	output, err := processor.Process(input)
	if err != nil {
		t.Fatalf("process failed: %v", err)
	}

	// Should remain JPEG effectively (or at least decodeable)
	_, format, err := image.Decode(output)
	if err != nil {
		t.Fatalf("failed to decode output: %v", err)
	}
	if format != "jpeg" {
		t.Errorf("expected jpeg, got %s", format)
	}
}

func TestImmutableFilters(t *testing.T) {
	processor := NewImageProcessor()

	// Branch 1
	p1 := processor.AddFilter(Resize(10, 10))

	// Branch 2
	p2 := processor.AddFilter(Resize(20, 20))

	if len(processor.filters) != 0 {
		t.Errorf("base processor should have 0 filters, got %d", len(processor.filters))
	}
	if len(p1.filters) != 1 {
		t.Errorf("p1 should have 1 filter, got %d", len(p1.filters))
	}
	if len(p2.filters) != 1 {
		t.Errorf("p2 should have 1 filter, got %d", len(p2.filters))
	}
}

func TestValidateRatio(t *testing.T) {
	// 1:1 Image
	squareImg, _ := createTestImage(100, 100, "png")
	// 16:9 Image (approx 100x56)
	landscapeImg, _ := createTestImage(160, 90, "png")

	processor := NewImageProcessor()

	// Test Valid 1:1
	pSquare := processor.AddFilter(ValidateRatio("1:1"))
	if _, err := pSquare.Process(squareImg); err != nil {
		t.Errorf("expected success for 1:1 check on square image, got %v", err)
	}

	// Use fresh reader for fail case
	squareImgFail, _ := createTestImage(100, 100, "png")
	// Test Invalid 16:9 check on Square Image
	pLandscape := processor.AddFilter(ValidateRatio("16:9"))
	if _, err := pLandscape.Process(squareImgFail); err == nil {
		t.Error("expected error for 16:9 check on square image, got nil")
	}

	// Test Valid 16:9
	if _, err := pLandscape.Process(landscapeImg); err != nil {
		t.Errorf("expected success for 16:9 check on landscape image, got %v", err)
	}
}

func TestVariadicAddFilter(t *testing.T) {
	processor := NewImageProcessor()
	p := processor.AddFilter(
		Resize(10, 10),
		FormatWebp(),
	)

	if len(p.filters) != 2 {
		t.Errorf("expected 2 filters, got %d", len(p.filters))
	}
}
