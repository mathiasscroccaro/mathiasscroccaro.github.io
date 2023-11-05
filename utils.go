package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

func imagetoBase64(imagePath string) string {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Fatalf("Image file %s does not exist", imagePath)
	}
	imageBytes, _ := os.ReadFile(imagePath)

	base64Image := base64.StdEncoding.EncodeToString(imageBytes)
	return string(base64Image)
}

func getAllHtmlImages(html []byte) []string {
	var images []string
	document, err := goquery.NewDocumentFromReader(bytes.NewReader(html))

	if err != nil {
		log.Fatal(err)
	}

	document.Find("img").Each(func(i int, s *goquery.Selection) {
		images = append(images, s.AttrOr("src", ""))
	})
	return images
}

func findImagesInRepositoryAndReplaceToBase64IntoHTML(html []byte) []byte {
	images := getAllHtmlImages(html)
	for _, image := range images {
		imagePath := filepath.Join("./posts/images", filepath.Base(image))
		base64Image := "data:image/jpeg;base64," + imagetoBase64(imagePath)
		html = bytes.Replace(html, []byte(image), []byte(base64Image), 1)
	}
	return html
}

func renderImage(w io.Writer, p *ast.Image, entering bool) {
	if entering {
		io.WriteString(w, fmt.Sprintf("<div class=\"image-container\"><img src=%s>", p.Destination))
	} else {
		io.WriteString(w, "</div>")
	}
}

func myRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if para, ok := node.(*ast.Image); ok {
		renderImage(w, para, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func newCustomizedRender() *html.Renderer {
	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: myRenderHook,
	}
	return html.NewRenderer(opts)
}

func mdToHTML(md []byte) []byte {
	renderer := newCustomizedRender()
	html := markdown.ToHTML(md, nil, renderer)

	return html
}

func getAllFilesInDirByPattern(dirPath, pattern string) []string {
	files, err := filepath.Glob(filepath.Join(dirPath, pattern))

	if err != nil || files == nil {
		log.Fatalf("It wasn't found any files with pattern %s in the directory %s", pattern, dirPath)
	}

	filesBaseNames := []string{}

	for _, file := range files {
		filesBaseNames = append(filesBaseNames, filepath.Base(file))
	}

	return filesBaseNames
}

func getAllMarkdownFilesInDir(dirPath string) []string {
	return getAllFilesInDirByPattern(dirPath, "*.md")
}

func getAllHtmlFilesInDir(dirPath string) []string {
	return getAllFilesInDirByPattern(dirPath, "*.html")
}

func createDirectory(dirPath string) string {
	os.MkdirAll(dirPath, os.ModePerm)
	return dirPath
}

func copyFile(originPath, destinationPath string) {
	file, err := os.ReadFile(originPath)

	if err != nil {
		log.Fatal(err)
	}

	createDirectory(filepath.Dir(destinationPath))
	os.WriteFile(destinationPath, file, 0644)
}

func copyStaticFiles(outputDirectory string) {
	copyFile("./assets/favicon.ico", filepath.Join(outputDirectory, "favicon.ico"))
	copyFile("./assets/styles.css", filepath.Join(outputDirectory, "styles.css"))
	copyFile("./assets/me.jpg", filepath.Join(outputDirectory, "me.jpg"))
}

func BuildStaticPortfolio(outputDirectory string) {
	createDirectory(outputDirectory)

	NewPortfolioTemplate().RenderPortfolioAndWriteToDir(outputDirectory)
	NewBlogPosts().RenderPosts(filepath.Join(outputDirectory, "posts"))

	copyStaticFiles(outputDirectory)
}
