package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	xhtml "golang.org/x/net/html"
)

var debug *log.Logger
var info *log.Logger
var warning *log.Logger
var errorLog *log.Logger

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
	debug = log.New(io.Discard, "DEBUG: ", log.LstdFlags|log.Lmsgprefix)
	info = log.New(io.Discard, "INFO: ", log.LstdFlags|log.Lmsgprefix)
	warning = log.New(os.Stderr, "WARNING: ", log.LstdFlags|log.Lmsgprefix)
	errorLog = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lmsgprefix)

	file, err := os.Open(os.Args[1])
	if err != nil {
		errorLog.Fatal(err)
	}
	doc, err := xhtml.Parse(file)
	if err != nil {
		errorLog.Fatal(err)
	}

	var rewriteRules map[string]string
	json.Unmarshal([]byte(rewrites), &rewriteRules)

	// Get and rewrite HTML tree.
	RewriteCrawler(doc, wd, rewriteRules)

	// Output or write updated HTML.
	var buffer bytes.Buffer
	writer := io.Writer(&buffer)
	xhtml.Render(writer, doc)
	fmt.Println(buffer.String())
}
