package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetPostNumber(t *testing.T) {
	post := Post{
		MarkdownPath: "./fixtures/posts/1.md",
	}
	if post.PostNumber() != 1 {
		t.Errorf("expected 1, got %d", post.PostNumber())
	}
}

func TestCreateBlogPosts(t *testing.T) {
	outputDir := "/tmp/posts"

	NewBlogPosts().RenderPosts(outputDir)

	files := getAllHtmlFilesInDir(outputDir)

	for i, file := range files {
		content, err := os.ReadFile(filepath.Join(outputDir, file))

		if err != nil || content == nil {
			t.Errorf("expected file %s to exist with content inside it", files[i])
		}

		os.Remove(filepath.Join(outputDir, files[i]))
	}
}
