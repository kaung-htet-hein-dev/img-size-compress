package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ImageStat struct {
	Name     string `json:"name"`
	Original int64  `json:"original"`
	Final    int64  `json:"final"`
}

func main() {
	targetDir := os.Args[1]
	var stats []ImageStat

	entries, err := os.ReadDir(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".jpg" && ext != ".png" {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		oldSize := info.Size()
		newSize := oldSize / 2 // Placeholder for actual optimization logic

		stats = append(stats, ImageStat{
			Name:     entry.Name(),
			Original: oldSize,
			Final:    newSize,
		})
	}

	if len(stats) == 0 {
		fmt.Fprint(os.Stderr, "Error: No images found in the specified directory.\n")
		os.Exit(1)
	}

	output, _ := json.Marshal(stats)
	fmt.Println(string(output))
}
