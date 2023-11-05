package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	outputDir := "./docs"

	BuildStaticPortfolio(outputDir)

	fs := http.FileServer(http.Dir(outputDir))
	http.Handle("/", fs)

	fmt.Print("Listening on http://localhost:3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
