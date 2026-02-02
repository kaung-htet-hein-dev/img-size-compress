package main

import (
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type ImageStat struct {
	Name     string `json:"name"`
	Original int64  `json:"original"`
	Final    int64  `json:"final"`
}

func compressImage(path string) (int64, int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to stat original file: %w", err)
	}
	originalSize := fileInfo.Size()

	src, err := imaging.Open(path)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to open image: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(path))

	tempFile, err := os.CreateTemp(filepath.Dir(path), "compress-*.tmp")
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath) // Clean up temp file

	switch ext {
	case ".jpg", ".jpeg":
		// Quality 75 provides a good balance between file size and visual quality
		err = jpeg.Encode(tempFile, src, &jpeg.Options{Quality: 75})
	case ".png":
		// PNG uses lossless compression; use maximum compression level
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		err = encoder.Encode(tempFile, src)
	}

	tempFile.Close()

	if err != nil {
		return 0, 0, fmt.Errorf("failed to encode image: %w", err)
	}

	compressedInfo, err := os.Stat(tempPath)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to stat compressed file: %w", err)
	}
	compressedSize := compressedInfo.Size()

	if compressedSize < originalSize {
		if err := copyFile(tempPath, path); err != nil {
			return 0, 0, fmt.Errorf("failed to replace original file: %w", err)
		}
		return originalSize, compressedSize, nil
	}

	return originalSize, originalSize, nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
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
