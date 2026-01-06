package images

import (
	"fmt"
	"math"

	"github.com/nfnt/resize"
)

// Resize resizes the image to the specified width and height.
// If one of width or height is 0, it preserves aspect ratio if the underlying library supports it (nfnt/resize does).
func Resize(width, height uint) Filter {
	return func(ctx *ProcessingContext) error {
		ctx.Image = resize.Resize(width, height, ctx.Image, resize.Lanczos3)
		return nil
	}
}

// ResizeAspectRatio resizes the image to fit within the specified width and height while maintaining aspect ratio.
// If keepRatio is false, it forces the dimensions to exact width and height.
// If keepRatio is true, it scales the image such that it fits within width x height, preserving the original aspect ratio.
func ResizeAspectRatio(width, height uint, keepRatio bool) Filter {
	return func(ctx *ProcessingContext) error {
		if !keepRatio {
			return Resize(width, height)(ctx)
		}

		// Calculate new dimensions while keeping aspect ratio
		ctx.Image = resize.Thumbnail(width, height, ctx.Image, resize.Lanczos3)
		return nil
	}
}

// FormatWebp converts the output format to WebP.
func FormatWebp() Filter {
	return func(ctx *ProcessingContext) error {
		ctx.OutputFormat = "webp"
		return nil
	}
}

// ValidateRatio validates if the image aspect ratio matches the expected ratio.
// Supported ratios: "1:1", "16:9", "9:16".
func ValidateRatio(ratio string) Filter {
	return func(ctx *ProcessingContext) error {
		bounds := ctx.Image.Bounds()
		width := float64(bounds.Dx())
		height := float64(bounds.Dy())

		var expectedRatio float64
		switch ratio {
		case "1:1":
			expectedRatio = 1.0
		case "16:9":
			expectedRatio = 16.0 / 9.0
		case "9:16":
			expectedRatio = 9.0 / 16.0
		default:
			return fmt.Errorf("unsupported ratio format: %s", ratio)
		}

		actualRatio := width / height

		// Use a nice epsilon for float comparison
		epsilon := 0.01
		if math.Abs(actualRatio-expectedRatio) > epsilon {
			return fmt.Errorf("image ratio mismatch: expected %s, got %.2f", ratio, actualRatio)
		}

		return nil
	}
}
