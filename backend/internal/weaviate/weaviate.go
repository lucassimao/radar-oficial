package weaviate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

func UploadDir(dir, description, entity string) error {

	cfg := weaviate.Config{
		Host:   os.Getenv("WEAVIATE_HOST"),
		Scheme: "http",
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create weaviate client: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Error reading directory: %v", err)

	}

	var objects []*models.Object
	// Regular expression to match the pattern "page_X.md"
	// This will match files like page_1.md, page_01.md, page_001.md, etc.
	pagePattern := regexp.MustCompile(`^page_(\d+)\.md$`)

	for _, entry := range entries {
		_, err := entry.Info()
		if err != nil {
			return fmt.Errorf("Error getting file info: %v", err)
		}

		// Check if the file matches our pattern
		matches := pagePattern.FindStringSubmatch(entry.Name())
		if matches == nil || len(matches) < 2 {
			continue // Skip files that don't match the pattern
		}

		// Extract the page number
		pageNumStr := matches[1]
		pageNum, err := strconv.Atoi(pageNumStr)
		if err != nil {
			return fmt.Errorf("Error parsing page number from %s: %v", entry.Name(), err)
		}

		// Get full path to the file
		filePath := filepath.Join(dir, entry.Name())

		// Read file content as string
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("Error reading file %s: %v", entry.Name(), err)
		}

		object := &models.Object{
			Class: "diarios",
			Properties: map[string]any{
				"description": description,
				"entity":      entity,
				"content":     string(content),
				"page":        pageNum,
			},
		}

		objects = append(objects, object)
	}

	// batch write items
	batchRes, err := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
	if err != nil {
		return fmt.Errorf("failed to batch upload objects: %v", err)
	}
	for _, res := range batchRes {
		if res.Result.Errors != nil {
			message := res.Result.Errors.Error[0].Message
			return fmt.Errorf("failed to batch upload objects: %s", message)
		}
	}

	return nil
}
