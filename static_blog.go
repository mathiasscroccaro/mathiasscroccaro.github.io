package main

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/ast"
	"log"
	"os"
	"io"
	"path/filepath"
	"html/template"
	"bytes"
	"strconv"
	"fmt"
	"encoding/base64"
	"github.com/PuerkitoBio/goquery"
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

func ImagetoBase64(imagePath string) string {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Fatalf("Image file %s does not exist", imagePath)
	}
	imageBytes, _ := os.ReadFile(imagePath)

	base64Image := base64.StdEncoding.EncodeToString(imageBytes)
	return string(base64Image)	
}

func GetAllHtmlImages(html []byte) []string {
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

func FindImagesInRepositoryAndReplaceToBase64IntoHTML(html []byte) []byte {
	images := GetAllHtmlImages(html)
	for _, image := range images {
		imagePath := filepath.Join("./posts/images", filepath.Base(image))
		base64Image := "data:image/jpeg;base64," + ImagetoBase64(imagePath)
		html = bytes.Replace(html, []byte(image), []byte(base64Image), 1)
	}
	return html
}

func (p Post) CreateHTMLFile() {
	createDirectory(filepath.Dir(p.HTMLPath))
	html := p.MarkdownToHTML()
	html = append([]byte("<article>"), html...)
	html = append(html, []byte("</article>")...)

	if p.PostNumber() != 1 {
		html = append(html, []byte(fmt.Sprintf("\n" + `<span hx-trigger="revealed" hx-get="/posts/%d.html" hx-swap="afterend"></span>`, p.PostNumber() - 1))...)
	}
	html = FindImagesInRepositoryAndReplaceToBase64IntoHTML(html)
	
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

type PortfolioTemplate struct {
	RootTemplatePath string
	HomepageTemplatePath string
	AboutTemplatePath string
	ProjectsTemplatePath string
}

func NewPortfolioTemplate() PortfolioTemplate {
	return PortfolioTemplate{
		RootTemplatePath: "./templates/root.tmpl",
		HomepageTemplatePath: "./templates/home.html",
		AboutTemplatePath: "./templates/about.html",
		ProjectsTemplatePath: "./templates/projects.tmpl",
	}
}

func (bt PortfolioTemplate) RenderHomePage() []byte {
	return bt.RenderPage(bt.HomepageTemplatePath)
}

func (bt PortfolioTemplate) RenderAboutPage() []byte {
	return bt.RenderPage(bt.AboutTemplatePath)
}

func (bt PortfolioTemplate) RenderProjectsPage() []byte {
	var buffer bytes.Buffer
	
	templateText := string(bt.RenderPage(bt.ProjectsTemplatePath))

	tmpl, err := template.New("projectsTemplate").Parse(templateText)
	if err != nil {
		log.Fatal(err)
	}
	
	files := getAllFilesInDirByPattern("./posts", "*.md")
	lastPostNumber := strconv.Itoa(len(files))

	tmpl.Execute(&buffer, lastPostNumber)

	return buffer.Bytes()
}

func (bt PortfolioTemplate) RenderPage(pageTemplatePath string) []byte {
	var buffer bytes.Buffer

	tmpl, err := template.ParseFiles(bt.RootTemplatePath)

	if err != nil {
		log.Fatal(err)
	}

	pageTemplate, err := os.ReadFile(pageTemplatePath)

	if err != nil {
		log.Fatal(err)
	}

	tmpl.Execute(&buffer, template.HTML(string(pageTemplate)))

	return buffer.Bytes()
}

func (bt PortfolioTemplate) RenderPortfolioAndWriteToDir(outputDir string) {
	createDirectory(outputDir)
	os.WriteFile(filepath.Join(outputDir, "index.html"), bt.RenderHomePage(), 0644)
	os.WriteFile(filepath.Join(outputDir, "about.html"), bt.RenderAboutPage(), 0644)
	os.WriteFile(filepath.Join(outputDir, "projects.html"), bt.RenderProjectsPage(), 0644)
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

func GetMarkdownFileTransformedToHTML(mdPath string) string {
	mdFile, _ := os.ReadFile(mdPath)
	html := mdToHTML(mdFile)
	return string(html)
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

func GetAllMarkdownFilesInDir(dirPath string) []string {
	return getAllFilesInDirByPattern(dirPath, "*.md")
}

func GetAllHtmlFilesInDir(dirPath string) []string {
	return getAllFilesInDirByPattern(dirPath, "*.html")
}

func createDirectory(dirPath string) string {
	os.MkdirAll(dirPath, os.ModePerm)
	return dirPath
}

func readMarkdownDirectoryAndCreatePosts(markdownDir, outputDirectory string) {
	for _, markdownFile := range GetAllMarkdownFilesInDir(markdownDir) {

		htmlFile := markdownFile[:len(markdownFile)-3] + ".html"
		html := GetMarkdownFileTransformedToHTML(filepath.Join(markdownDir, markdownFile))
		os.WriteFile(filepath.Join(outputDirectory, htmlFile), []byte(html), 0644)
	}
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
	// copyFile("./assets/favicon.ico", filepath.Join(outputDirectory, "favicon.ico"))
	copyFile("./assets/styles.css", filepath.Join(outputDirectory, "styles.css"))
	copyFile("./assets/me.jpg", filepath.Join(outputDirectory, "me.jpg"))
}

func BuildStaticPortfolio(outputDirectory string) {
	createDirectory(outputDirectory)

	NewPortfolioTemplate().RenderPortfolioAndWriteToDir(outputDirectory)
	NewBlogPosts().RenderPosts(filepath.Join(outputDirectory, "posts"))

	copyStaticFiles(outputDirectory)
}
