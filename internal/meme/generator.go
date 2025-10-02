package meme

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Generator handles meme generation
type Generator struct {
	inputPath string
	outputDir string
}

// NewGenerator creates a new meme generator
func NewGenerator(inputPath, outputDir string) *Generator {
	return &Generator{
		inputPath: inputPath,
		outputDir: outputDir,
	}
}

// GenerateMemes processes all images in the input folder and generates memes
func (g *Generator) GenerateMemes(textTop, textBottom string) error {
	// Walk through the input directory
	return filepath.WalkDir(g.inputPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Check if file is an image
		if !g.isImageFile(path) {
			return nil
		}

		// Generate meme for this image
		return g.generateMemeFromFile(path, textTop, textBottom)
	})
}

// isImageFile checks if a file is an image based on its extension
func (g *Generator) isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

// generateMemeFromFile creates a meme from an image file
func (g *Generator) generateMemeFromFile(imagePath, textTop, textBottom string) error {
	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	// Decode the image
	var img image.Image
	ext := strings.ToLower(filepath.Ext(imagePath))

	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return fmt.Errorf("unsupported image format: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// Create a new RGBA image
	bounds := img.Bounds()
	m := image.NewRGBA(bounds)
	draw.Draw(m, bounds, img, bounds.Min, draw.Src)

	// Add text to the image
	if textTop != "" {
		g.addText(m, textTop, bounds.Dx()/2, 50, bounds.Dx(), bounds.Dy())
	}

	if textBottom != "" {
		g.addText(m, textBottom, bounds.Dx()/2, bounds.Dy()-50, bounds.Dx(), bounds.Dy())
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Save the generated meme
	outputPath := filepath.Join(g.outputDir, filepath.Base(imagePath))
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Encode and save the image
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(outputFile, m, &jpeg.Options{Quality: 90})
	case ".png":
		err = png.Encode(outputFile, m)
	}

	if err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	fmt.Printf("Generated meme: %s\n", outputPath)
	return nil
}

// addText adds text to an image with proportional sizing
func (g *Generator) addText(img *image.RGBA, text string, x, y, imgWidth, imgHeight int) {
	// Calculate font size as 7% of the image height
	fontSize := imgHeight * 7 / 100
	if fontSize < 10 {
		fontSize = 10 // Minimum font size
	}

	// Calculate text position with better spacing
	// Position text at 7% from top/bottom edges
	yPos := y
	if y < imgHeight/2 {
		// Top text - position at 7% from top
		yPos = imgHeight * 7 / 100
	} else {
		// Bottom text - position at 7% from bottom
		yPos = imgHeight - (imgHeight * 7 / 100)
	}

	// Calculate text width approximation (7 pixels per character for basicfont.Face7x13)
	textWidth := len(text) * 7
	textHeight := 13 // Height of basicfont.Face7x13

	// Draw semi-transparent white background
	// Create a slightly larger background rectangle
	bgX1 := x - textWidth/2 - 10
	bgY1 := yPos - textHeight - 5
	bgX2 := x + textWidth/2 + 10
	bgY2 := yPos + 5

	// Ensure background stays within image bounds
	if bgX1 < 0 {
		bgX1 = 0
	}
	if bgY1 < 0 {
		bgY1 = 0
	}
	if bgX2 > imgWidth {
		bgX2 = imgWidth
	}
	if bgY2 > imgHeight {
		bgY2 = imgHeight
	}

	// Draw semi-transparent white background
	for i := bgX1; i < bgX2; i++ {
		for j := bgY1; j < bgY2; j++ {
			// Get existing pixel
			existing := img.RGBAAt(i, j)
			// Blend with semi-transparent white (128 alpha)
			alpha := 128
			blend := color.RGBA{
				R: uint8((int(existing.R)*(255-alpha) + 255*alpha) / 255),
				G: uint8((int(existing.G)*(255-alpha) + 255*alpha) / 255),
				B: uint8((int(existing.B)*(255-alpha) + 255*alpha) / 255),
				A: existing.A,
			}
			img.SetRGBA(i, j, blend)
		}
	}

	// Draw multiple layers of black outline for better visibility
	// Use a 3x3 grid for a thicker outline
	// Increase outline thickness for larger text appearance
	outlineThickness := fontSize / 20
	if outlineThickness < 2 {
		outlineThickness = 2
	}

	for i := -outlineThickness; i <= outlineThickness; i++ {
		for j := -outlineThickness; j <= outlineThickness; j++ {
			// Skip the center and inner positions to create a hollow outline effect
			if (i >= -1 && i <= 1) && (j >= -1 && j <= 1) {
				continue
			}
			blackDrawer := &font.Drawer{
				Dst:  img,
				Src:  image.NewUniform(color.RGBA{0, 0, 0, 255}), // Black outline
				Face: basicfont.Face7x13,
				Dot:  fixed.Point26_6{X: fixed.I(x - (len(text)*7)/2 + i), Y: fixed.I(yPos + j)},
			}
			blackDrawer.DrawString(text)
		}
	}

	// Draw additional layer closer to text for a cleaner outline
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue // Skip center position for outline
			}
			blackDrawer := &font.Drawer{
				Dst:  img,
				Src:  image.NewUniform(color.RGBA{0, 0, 0, 255}), // Black outline
				Face: basicfont.Face7x13,
				Dot:  fixed.Point26_6{X: fixed.I(x - (len(text)*7)/2 + i), Y: fixed.I(yPos + j)},
			}
			blackDrawer.DrawString(text)
		}
	}

	// Draw the main white text in the center
	// Draw multiple times to make it appear larger
	charWidth := fontSize / 10
	if charWidth < 1 {
		charWidth = 1
	}

	for i := -charWidth; i <= charWidth; i++ {
		for j := -charWidth; j <= charWidth; j++ {
			whiteDrawer := &font.Drawer{
				Dst:  img,
				Src:  image.NewUniform(color.RGBA{255, 255, 255, 255}), // White text
				Face: basicfont.Face7x13,
				Dot:  fixed.Point26_6{X: fixed.I(x - (len(text)*7)/2 + i), Y: fixed.I(yPos + j)},
			}
			whiteDrawer.DrawString(text)
		}
	}
}

// CreateMemeImage creates a meme image with the given text and saves it to the specified path
func CreateMemeImage(textTop, textBottom, outputPath string) error {
	// Create a default template
	width, height := 800, 600
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with a light gray background
	bgColor := color.RGBA{200, 200, 200, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Add border
	borderColor := color.RGBA{0, 0, 0, 255}
	for i := 0; i < width; i++ {
		img.Set(i, 0, borderColor)
		img.Set(i, height-1, borderColor)
	}
	for i := 0; i < height; i++ {
		img.Set(0, i, borderColor)
		img.Set(width-1, i, borderColor)
	}

	// Add text to the image
	generator := &Generator{}
	if textTop != "" {
		generator.addText(img, textTop, width/2, 50, width, height)
	}

	if textBottom != "" {
		generator.addText(img, textBottom, width/2, height-50, width, height)
	}

	// Create output directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Save the generated meme
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Encode and save the image as PNG
	if err := png.Encode(outputFile, img); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	return nil
}

// CreateMemeFromTemplate creates a meme using a template image
func CreateMemeFromTemplate(templateName, textTop, textBottom, outputPath string) error {
	// Load the template image
	templateImg, err := LoadTemplateImage(templateName)
	if err != nil {
		return fmt.Errorf("failed to load template image: %w", err)
	}

	// Create a new RGBA image from the template
	bounds := templateImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, templateImg, bounds.Min, draw.Src)

	// Add text to the image
	generator := &Generator{}
	if textTop != "" {
		generator.addText(img, textTop, bounds.Dx()/2, 50, bounds.Dx(), bounds.Dy())
	}

	if textBottom != "" {
		generator.addText(img, textBottom, bounds.Dx()/2, bounds.Dy()-50, bounds.Dx(), bounds.Dy())
	}

	// Create output directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Save the generated meme
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Encode and save the image as PNG
	if err := png.Encode(outputFile, img); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	return nil
}

// LoadTemplateImage loads a template image by name
func LoadTemplateImage(templateName string) (image.Image, error) {
	// Construct the path to the template images directory
	imagesDir := fmt.Sprintf("./data/templates/%s/images", templateName)

	// Check if images directory exists
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("images directory not found for template %s", templateName)
	}

	// Walk the directory to find the first image file
	var imagePath string
	err := filepath.WalkDir(imagesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Check if file is an image
		if isImageFile(path) {
			imagePath = path
			return filepath.SkipDir // Stop walking after finding the first image
		}

		return nil
	})

	if err != nil || imagePath == "" {
		return nil, fmt.Errorf("no image found for template %s", templateName)
	}

	// Load and return the image
	return loadSingleImage(imagePath)
}

// loadSingleImage loads a single image file
func loadSingleImage(imagePath string) (image.Image, error) {
	// Check if file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil, err
	}

	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	// Decode the image
	ext := strings.ToLower(filepath.Ext(imagePath))

	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Decode(file)
	case ".png":
		return png.Decode(file)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}
}

// isImageFile checks if a file is an image based on its extension
func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}
