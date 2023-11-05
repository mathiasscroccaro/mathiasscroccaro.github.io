package main

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type PortfolioTemplate struct {
	RootTemplatePath     string
	HomepageTemplatePath string
	AboutTemplatePath    string
	ProjectsTemplatePath string
}

func NewPortfolioTemplate() PortfolioTemplate {
	return PortfolioTemplate{
		RootTemplatePath:     "./templates/root.tmpl",
		HomepageTemplatePath: "./templates/home.html",
		AboutTemplatePath:    "./templates/about.html",
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
