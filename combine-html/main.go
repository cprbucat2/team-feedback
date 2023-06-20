package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"os"
	"path"

	xhtml "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var debug *log.Logger
var info *log.Logger
var warning *log.Logger
var errorLog *log.Logger

func newTextNode(text string) *xhtml.Node {
	node := xhtml.Node{
		Type: xhtml.TextNode,
		Data: text,
	}
	return &node
}

func rewriteCSSLink(node *xhtml.Node, file string) error {
	var text string

	if fh, err := os.Open(file); err == nil {
		if buf, err := io.ReadAll(fh); err == nil {
			text = string(buf)
		} else {
			warning.Println(err)
			return err
		}
	} else {
		warning.Println(err)
		return err
	}

	newChild := newTextNode(text)
	node.DataAtom = atom.Style
	node.Data = "style"
	newAttrs := []xhtml.Attribute{}
	for _, at := range node.Attr {
		if at.Key != "rel" && at.Key != "href" {
			newAttrs = append(newAttrs, at)
		}
	}
	node.Attr = newAttrs
	node.AppendChild(newChild)

	return nil
}

func rewriteScript() {}

func imageText(file string) string {
	if path.Ext(file) == ".svg" {
		if handle, err := os.Open(file); err != nil {
			warning.Println(err)
		} else {
			if buf, err := io.ReadAll(handle); err != nil {
				warning.Println(err)
			} else {
				return "data:image/svg+xml," + html.EscapeString(string(buf))
			}
		}
	} else if path.Ext(file) == ".png" {
		if handle, err := os.Open(file); err != nil {
			warning.Println(err)
		} else {
			if buf, err := io.ReadAll(handle); err != nil {
				warning.Println(err)
			} else {
				return "data:image/png;base64," + html.EscapeString(base64.URLEncoding.EncodeToString(buf))
			}
		}
	} else {
		warning.Println("Unknown image extension for file ", file)
	}
	return ""
}

func rewriteFaviconLink(node *xhtml.Node, file string) {
	for i := range node.Attr {
		if node.Attr[i].Key == "href" {
			node.Attr[i].Val = imageText(file)
		} else if node.Attr[i].Key == "type" {
			if path.Ext(file) == ".svg" {
				node.Attr[i].Val = "image/svg+xml"
			} else if path.Ext(file) == ".png" {
				node.Attr[i].Val = "image/png"
			}
		}
	}
}

func rewriteIMG(node *xhtml.Node, file string) {
	for i := range node.Attr {
		if node.Attr[i].Key == "src" {
			node.Attr[i].Val = imageText(file)
		}
	}
}

func rewritePath(rewriteRules map[string]string, old string) string {
	if href, ok := rewriteRules[old]; ok {
		info.Printf("Rewriting %s -> %s.\n", old, href)
		return href
	} else {
		return old
	}
}

func crawler(n *xhtml.Node, stem string, rewriteRules map[string]string) bool {
	if n.Type == xhtml.ElementNode && n.DataAtom == atom.Link {
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

			if newHref := rewritePath(rewriteRules, href); newHref != "" {
				filePath := path.Join(stem, newHref)
				debug.Printf("Calculated pathname: %v\n", filePath)
				rewriteCSSLink(n, filePath)
			} else {
				return false
			}
		} else if rel == "shortcut icon" {
			info.Printf("Found a favicon <link> tag with href='%s'.\n", href)

			if newHref := rewritePath(rewriteRules, href); newHref != "" {
				filePath := path.Join(stem, newHref)
				debug.Printf("Calculated pathname: %v\n", filePath)
				rewriteFaviconLink(n, filePath)
			} else {
				return false
			}
		}
		return true
	}

	for child := n.FirstChild; child != nil; {
		if !crawler(child, stem, rewriteRules) {
			info.Printf("Removing <%s> element.\n", child.Data)
			child = child.NextSibling
			if child != nil {
				n.RemoveChild(child.PrevSibling)
			} else {
				n.RemoveChild(n.LastChild)
			}
		} else {
			child = child.NextSibling
		}
	}
	return true
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
	doc, err := xhtml.Parse(file)
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
	xhtml.Render(writer, doc)
	fmt.Println(buffer.String())
}
