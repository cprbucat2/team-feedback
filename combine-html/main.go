package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/net/html"
)

func main() {
	// TODO: Process argument to write in place.
	// TODO: Process argument comtaining source rewrites.
	if len(os.Args) < 1 {
		fmt.Printf("Usage: %s index.html\n", os.Args[0])
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	doc, err := html.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	// Get a list of internal sources.

	// Filter or rewrite internal sources.

	// Replace source links with immediate versions.

	// Output or write updated HTML.
	var buffer bytes.Buffer
	writer := io.Writer(&buffer)
	html.Render(writer, doc)
	fmt.Print(buffer.String())
}
