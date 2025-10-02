package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"memes-generator/internal/domain"
	"memes-generator/internal/meme"
)

func main() {
	var memePath string

	flag.StringVar(&memePath, "meme-path", "", "Path to meme directory")
	flag.Parse()

	if memePath == "" {
		fmt.Println("Usage: generate-meme --meme-path <path-to-meme-directory>")
		fmt.Println("Example: ./generate-meme --meme-path data/memes/meme_1759442111813095000")
		os.Exit(1)
	}

	// Check if meme path exists
	if _, err := os.Stat(memePath); os.IsNotExist(err) {
		log.Fatalf("Meme path does not exist: %s", memePath)
	}

	// Read metadata.json from the meme directory
	metadataPath := filepath.Join(memePath, "metadata.json")
	metadataFile, err := os.Open(metadataPath)
	if err != nil {
		log.Fatalf("Failed to open metadata file: %v", err)
	}
	defer metadataFile.Close()

	var memeEntity domain.Meme
	decoder := json.NewDecoder(metadataFile)
	if err := decoder.Decode(&memeEntity); err != nil {
		log.Fatalf("Failed to decode metadata: %v", err)
	}

	// Use the same output directory (overwrite existing images)
	imageDir := filepath.Join(memePath, "images")
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Check if there are any source images
	imagesDir := filepath.Join(memePath, "images")
	sourceImagesFound := false

	if _, err := os.Stat(imagesDir); err == nil {
		entries, err := os.ReadDir(imagesDir)
		if err == nil {
			for _, entry := range entries {
				// Skip generated_meme.png as it's an output file, not a source
				if !entry.IsDir() && isImageFile(entry.Name()) && entry.Name() != "generated_meme.png" {
					sourceImagesFound = true
					break
				}
			}
		}
	}

	if sourceImagesFound {
		// Initialize meme generator with existing images
		generator := meme.NewGenerator(imagesDir, imageDir)

		// Generate memes using text from metadata
		fmt.Printf("Regenerating memes for ID: %s\n", memeEntity.ID)
		fmt.Printf("Using text: Top='%s', Bottom='%s'\n", memeEntity.TextTop, memeEntity.TextBottom)
		fmt.Printf("Input path: %s\n", imagesDir)
		fmt.Printf("Output path: %s\n", imageDir)

		if err := generator.GenerateMemes(memeEntity.TextTop, memeEntity.TextBottom); err != nil {
			log.Fatalf("Failed to generate memes: %v", err)
		}
	} else {
		// No source images found, create a meme from template or scratch
		fmt.Printf("No source images found for ID: %s\n", memeEntity.ID)
		fmt.Printf("Creating meme using template '%s' with text: Top='%s', Bottom='%s'\n", memeEntity.Template, memeEntity.TextTop, memeEntity.TextBottom)
		fmt.Printf("Output path: %s\n", imageDir)

		// Try to create meme from template
		outputPath := filepath.Join(imageDir, "generated_meme.png")
		if err := meme.CreateMemeFromTemplate(memeEntity.Template, memeEntity.TextTop, memeEntity.TextBottom, outputPath); err != nil {
			// If template not found, create a simple default template
			fmt.Printf("Template '%s' not found, creating meme from scratch\n", memeEntity.Template)
			if err := meme.CreateMemeImage(memeEntity.TextTop, memeEntity.TextBottom, outputPath); err != nil {
				log.Fatalf("Failed to create meme from scratch: %v", err)
			}
		}

		fmt.Printf("Generated meme: %s\n", outputPath)
	}

	fmt.Printf("Memes regenerated successfully for ID: %s\n", memeEntity.ID)

	// Print information about generated files
	files, err := os.ReadDir(imageDir)
	if err == nil {
		fmt.Printf("Generated %d memes:\n", len(files))
		for _, file := range files {
			if !file.IsDir() && isImageFile(file.Name()) {
				fmt.Printf("  - %s\n", file.Name())
			}
		}
	}
}

// isImageFile checks if a file is an image based on its extension
func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}
