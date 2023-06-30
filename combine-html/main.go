package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	xhtml "golang.org/x/net/html"
)

func main() {
	// TODO: Process argument to write in place.
	// TODO: Process argument comtaining source rewrites.
	// TODO: Add option to specify working directory other than index html stem.
	// TODO: Add -l, --list option to list external resources.
	// TODO: Add -s, --dry-run to check which external resources are missing.
	// TODO: Add log level option.
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s index.html\n", os.Args[0])
		return
	}
	rewrites := `{"favicon.ico": ""}`
	wd := path.Dir(os.Args[1])

	// Setup loggers.
	LogDefault()

	file, err := os.Open(os.Args[1])
	if err != nil {
		errorLog.Fatal(err)
	}
	doc, err := xhtml.Parse(file)
	if err != nil {
		errorLog.Fatal(err)
	}

	var rewriteRules map[string]string
	if err := json.Unmarshal([]byte(rewrites), &rewriteRules); err != nil {
		errorLog.Fatal(err)
	}

	// Get and rewrite HTML tree.
	RewriteCrawler(doc, wd, rewriteRules)

	// Output or write updated HTML.
	var buffer bytes.Buffer
	writer := io.Writer(&buffer)
	if err := xhtml.Render(writer, doc); err != nil {
		errorLog.Fatal(err)
	}
	fmt.Println(buffer.String())
}
