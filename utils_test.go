package main

import (
	"os"
	"testing"
)

func TestGetAllMarkdownFilesInDir(t *testing.T) {
	files := getAllMarkdownFilesInDir("./posts/")
	if len(files) == 0 {
		t.Errorf("expected files, got %d", len(files))
	}
}

func TestImageToBase64(t *testing.T) {
	imageBase64 := imagetoBase64("./fixtures/base64_image_example.png")
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

	images := getAllHtmlImages(html)
	if len(images) == 0 {
		t.Errorf("expected images, got %d", len(images))
	}
	if images[0] != "/assets/img/MarineGEO_logo.png" {
		t.Errorf("expected /assets/img/MarineGEO_logo.png, got %s", images[0])
	}
}
