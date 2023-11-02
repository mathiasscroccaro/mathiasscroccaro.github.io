package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetMarkdownFileTransformedToHTML(t *testing.T) {
	html := GetMarkdownFileTransformedToHTML("./fixtures/basic-convertion.md")

	if html == "" {
		t.Errorf("expected html to not be empty")
	}

	expectedHtmlBytes, _ := os.ReadFile("./fixtures/basic-convertion.html")
	expectedHtml := string(expectedHtmlBytes)

	if html != expectedHtml {
		t.Errorf("expected %s, got %s", expectedHtml, html)
	}
}

func TestGetAllMarkdownFilesInDir(t *testing.T) {
	files := GetAllMarkdownFilesInDir("./posts/")
	if len(files) == 0 {
		t.Errorf("expected files, got %d", len(files))
	}
}

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

	files := GetAllHtmlFilesInDir(outputDir)

	for i, file := range files {
		content, err := os.ReadFile(filepath.Join(outputDir, file))

		if err != nil || content == nil {
			t.Errorf("expected file %s to exist with content inside it", files[i])
		}

		os.Remove(filepath.Join(outputDir, files[i]))
	}
}

func TestRenderHomePage(t *testing.T) {
	renderedHomePage := NewPortfolioTemplate().RenderHomePage()
	if string(renderedHomePage) == "" {
		t.Errorf("expected renderedHomePage to not be empty")
	}
}

func TestCreatePortfolio(t *testing.T) {
	outputDir := "./output/"

	NewPortfolioTemplate().RenderPortfolioAndWriteToDir(outputDir)

	files := GetAllHtmlFilesInDir(outputDir)

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

		// os.Remove(filepath.Join(outputDir, files[i]))
	}
}

func TestImageToBase64(t *testing.T) {
	imageBase64 := ImagetoBase64("./fixtures/base64_image_example.png")
	if imageBase64 == "" {
		t.Errorf("expected imageBase64 to not be empty")
	}
	expectedBase64Image := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="
	if imageBase64 != expectedBase64Image {
		t.Errorf("expected %s, got %s", expectedBase64Image, imageBase64)
	}
}

func TestGetAllHtmlImages(t *testing.T) {
	html, err := os.ReadFile("./fixtures/basic-convertion.html")
	
	if err != nil {
		t.Errorf("expected html to not be empty")
	}

	images := GetAllHtmlImages(html)
	if len(images) == 0 {
		t.Errorf("expected images, got %d", len(images))
	}
	if images[0] != "/assets/img/MarineGEO_logo.png" {
		t.Errorf("expected /assets/img/MarineGEO_logo.png, got %s", images[0])
	}
}
