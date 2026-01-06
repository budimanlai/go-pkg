package images

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"github.com/chai2010/webp"
)

// ProcessingContext holds the state of the image during processing.
type ProcessingContext struct {
	Image        image.Image
	OutputFormat string // "jpeg", "png", "webp", etc.
	Quality      int    // Quality for jpeg/webp (0-100). Default usually 75-90 depending on encoder.
}

// Filter is a function that modifies the ProcessingContext.
type Filter func(*ProcessingContext) error

// ImageProcessor handles the image processing pipeline.
type ImageProcessor struct {
	filters []Filter
}

// NewImageProcessor creates a new instance of ImageProcessor.
func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		filters: make([]Filter, 0),
	}
}

// AddFilter adds filters to the processing pipeline and returns a new ImageProcessor instance.
// It does not modify the current instance.
func (p *ImageProcessor) AddFilter(filters ...Filter) *ImageProcessor {
	// Create a new slice with existing filters
	newFilters := make([]Filter, len(p.filters), len(p.filters)+len(filters))
	copy(newFilters, p.filters)

	// Append new filters
	newFilters = append(newFilters, filters...)

	return &ImageProcessor{
		filters: newFilters,
	}
}

// Process processes the input image through the filter chain and returns the result.
func (p *ImageProcessor) Process(input io.Reader) (io.Reader, error) {
	// 1. Decode the input image
	img, format, err := image.Decode(input)
	if err != nil {
		return nil, err
	}

	// 2. Initialize context
	ctx := &ProcessingContext{
		Image:        img,
		OutputFormat: format, // Default to input format if possible, or we will handle fallback
		Quality:      90,     // Default quality
	}

	// 3. Apply filters
	for _, filter := range p.filters {
		if err := filter(ctx); err != nil {
			return nil, err
		}
	}

	// 4. Encode result
	var buf bytes.Buffer
	err = p.encode(&buf, ctx)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

// encode encodes the image in the context to the writer based on OutputFormat.
func (p *ImageProcessor) encode(w io.Writer, ctx *ProcessingContext) error {
	format := strings.ToLower(ctx.OutputFormat)

	switch format {
	case "webp":
		return webp.Encode(w, ctx.Image, &webp.Options{Quality: float32(ctx.Quality)})
	case "png":
		return png.Encode(w, ctx.Image)
	case "jpeg", "jpg":
		return jpeg.Encode(w, ctx.Image, &jpeg.Options{Quality: ctx.Quality})
	default:
		// Fallback to PNG if format is unknown or not explicitly supported
		return png.Encode(w, ctx.Image)
	}
}
