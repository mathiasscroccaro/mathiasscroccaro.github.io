package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Post struct {
	MarkdownPath string
	HTMLPath     string
}

func (p Post) ReadMarkdownFile() []byte {
	if _, err := os.Stat(p.MarkdownPath); os.IsNotExist(err) {
		log.Fatalf("Markdown file %s does not exist", p.MarkdownPath)
	}
	mdFile, _ := os.ReadFile(p.MarkdownPath)
	return mdFile
}

func (p Post) ReadHTMLFile() []byte {
	htmlFile, _ := os.ReadFile(p.HTMLPath)
	return htmlFile
}

func (p Post) MarkdownToHTML() []byte {
	mdFile := p.ReadMarkdownFile()
	html := mdToHTML(mdFile)
	return html
}

func (p Post) CreateHTMLFile() {
	createDirectory(filepath.Dir(p.HTMLPath))
	html := p.MarkdownToHTML()
	html = append([]byte("<article>"), html...)
	html = append(html, []byte("</article>")...)

	if p.PostNumber() != 1 {
		html = append(html, []byte(fmt.Sprintf("\n"+`<span hx-trigger="revealed" hx-get="/posts/%d.html" hx-swap="afterend"></span>`, p.PostNumber()-1))...)
	}
	html = findImagesInRepositoryAndReplaceToBase64IntoHTML(html)

	os.WriteFile(p.HTMLPath, html, 0644)
}

func (p Post) PostNumber() int {
	fileName := filepath.Base(p.MarkdownPath)
	value, err := strconv.Atoi(fileName[:len(fileName)-3])
	if err != nil {
		log.Fatalf("Impossible to convert the markdown file %s to an integer", fileName)
	}
	return value
}

type BlogPosts struct {
	Posts []Post
}

func NewBlogPosts() BlogPosts {
	return BlogPosts{}
}

func (bp BlogPosts) RenderPosts(htmlDir string) ([]Post, error) {
	posts := []Post{}
	markdownDir := "./posts"
	markdownFiles := getAllFilesInDirByPattern(markdownDir, "*.md")
	for _, markdownFile := range markdownFiles {
		fileName := markdownFile[:len(markdownFile)-3]

		post := Post{
			MarkdownPath: filepath.Join(markdownDir, fileName+".md"),
			HTMLPath:     filepath.Join(htmlDir, fileName+".html"),
		}

		post.CreateHTMLFile()
	}
	return posts, nil
}
