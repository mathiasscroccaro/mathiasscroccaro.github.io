package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRenderHomePage(t *testing.T) {
	renderedHomePage := NewPortfolioTemplate().RenderHomePage()
	if string(renderedHomePage) == "" {
		t.Errorf("expected renderedHomePage to not be empty")
	}
}

func TestCreatePortfolio(t *testing.T) {
	outputDir := "./test_output/"

	NewPortfolioTemplate().RenderPortfolioAndWriteToDir(outputDir)

	files := getAllHtmlFilesInDir(outputDir)

	expectedFiles := []string{
		"about.html",
		"index.html",
		"projects.html",
	}

	for i, expectedFile := range expectedFiles {
		if files[i] != expectedFile {
			t.Errorf("expected %s, got %s", expectedFile, files[i])
		}

		content, err := os.ReadFile(filepath.Join(outputDir, files[i]))

		if err != nil || content == nil {
			t.Errorf("expected file %s to exist with content inside it", files[i])
		}

		os.Remove(filepath.Join(outputDir, files[i]))
	}
}
