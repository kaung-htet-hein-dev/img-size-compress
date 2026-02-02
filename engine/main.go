package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/bimg"
)

type ImageStat struct {
	Name     string `json:"name"`
	Original int64  `json:"original"`
	Final    int64  `json:"final"`
}

func compressImage(path string) (int64, int64, error) {
	// Get original file size
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to stat original file: %w", err)
	}
	originalSize := fileInfo.Size()

	// Read image file
	buffer, err := bimg.Read(path)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read image: %w", err)
	}

	// Configure compression options
	options := bimg.Options{
		Quality:     75, // JPEG quality (1-100)
		Compression: 9,  // PNG compression level (0-9, 9 = best compression)
	}

	// Process image with bimg
	newImage, err := bimg.NewImage(buffer).Process(options)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to process image: %w", err)
	}

	compressedSize := int64(len(newImage))

	// Only replace if compressed version is smaller
	if compressedSize < originalSize {
		if err := os.WriteFile(path, newImage, 0644); err != nil {
			return 0, 0, fmt.Errorf("failed to write compressed image: %w", err)
		}
		return originalSize, compressedSize, nil
	}

	return originalSize, originalSize, nil
}

func main() {
	targetDir := os.Args[1]
	var compressionStats []ImageStat

	dirEntries, err := os.ReadDir(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
		os.Exit(1)
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		extension := strings.ToLower(filepath.Ext(entry.Name()))
		if extension != ".jpg" && extension != ".jpeg" && extension != ".png" {
			continue
		}

		filePath := filepath.Join(targetDir, entry.Name())
		originalSize, compressedSize, err := compressImage(filePath)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Skipping %s: %v\n", entry.Name(), err)
			continue
		}

		compressionStats = append(compressionStats, ImageStat{
			Name:     entry.Name(),
			Original: originalSize,
			Final:    compressedSize,
		})
	}

	if len(compressionStats) == 0 {
		fmt.Fprint(os.Stderr, "Error: No images found in the specified directory.\n")
		os.Exit(1)
	}

	jsonOutput, err := json.Marshal(compressionStats)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonOutput))
}
