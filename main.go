package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func initServer(outputDir string) {
	fs := http.FileServer(http.Dir(outputDir))
	http.Handle("/", fs)

	fmt.Println("Listening on http://localhost:3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var outputDir string
	var server bool

	flag.StringVar(&outputDir, "o", "./docs", "output directory")
	flag.BoolVar(&server, "s", false, "run server")

	flag.Parse()

	fmt.Println("Outputing to " + outputDir)
	BuildStaticPortfolio(outputDir)

	if server {
		initServer(outputDir)
	}
}
