package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var debug *log.Logger
var info *log.Logger
var warning *log.Logger
var errorLog *log.Logger

func newTextNode(text string) *html.Node {
	node := html.Node{
		Type: html.TextNode,
		Data: text,
	}
	return &node
}

func rewriteCSSLink(node *html.Node, file string) error {
	var text string

	if fh, err := os.Open(file); err == nil {
		if buf, err := io.ReadAll(fh); err == nil {
			text = string(buf)
		} else {
			warning.Print(err)
			return err
		}
	} else {
		warning.Print(err)
		return err
	}

	newChild := newTextNode(text)
	node.DataAtom = atom.Style
	node.Data = "style"
	newAttrs := []html.Attribute{}
	for _, at := range node.Attr {
		if at.Key != "rel" && at.Key != "href" {
			newAttrs = append(newAttrs, at)
		}
	}
	node.Attr = newAttrs
	node.AppendChild(newChild)

	return nil
}

func rewriteFaviconLink() {}
func rewriteIMG()         {}
func rewriteScript()      {}

func crawler(n *html.Node, stem string, rewriteRules map[string]string) {
	if n.Type == html.ElementNode && n.DataAtom == atom.Link {
		var rel string
		var href string
		for _, attr := range n.Attr {
			switch attr.Key {
			case "rel":
				rel = attr.Val
			case "href":
				href = attr.Val
			}
		}
		if rel == "stylesheet" {
			info.Printf("Found a stylesheet <link> tag with href='%s'.\n", href)
			if _, ok := rewriteRules[href]; ok {
				info.Printf("Rewriting %s -> %s.\n", href, rewriteRules[href])
				href = rewriteRules[href]
			}

			filePath := path.Join(stem, href)

			debug.Printf("Calculated pathname: %v\n", filePath)
			rewriteCSSLink(n, filePath)

		} else if rel == "shortcut icon" {
			info.Printf("Found a favicon <link> tag with href='%s'.\n", href)
			if _, ok := rewriteRules[href]; ok {
				info.Printf("Rewriting %s -> %s.\n", href, rewriteRules[href])
			}

			filePath := path.Join(stem, href)

			debug.Printf("Calculated pathname: %v\n", filePath)
		}
		return
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		crawler(child, stem, rewriteRules)
	}
}

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

	debug = log.New(io.Discard, "DEBUG: ", log.LstdFlags|log.Lmsgprefix)
	info = log.New(io.Discard, "INFO: ", log.LstdFlags|log.Lmsgprefix)
	warning = log.New(os.Stderr, "WARNING: ", log.LstdFlags|log.Lmsgprefix)
	errorLog = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lmsgprefix)

	file, err := os.Open(os.Args[1])
	if err != nil {
		errorLog.Fatal(err)
	}
	doc, err := html.Parse(file)
	if err != nil {
		errorLog.Fatal(err)
	}

	var rewriteRules map[string]string
	json.Unmarshal([]byte(rewrites), &rewriteRules)

	// Get a list of internal sources.
	crawler(doc, wd, rewriteRules)

	// Filter or rewrite internal sources.

	// Replace source links with immediate versions.

	// Output or write updated HTML.
	var buffer bytes.Buffer
	writer := io.Writer(&buffer)
	html.Render(writer, doc)
	fmt.Println(buffer.String())
}
