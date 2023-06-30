package main

import (
	"encoding/base64"
	"errors"
	"html"
	"io"
	"log"
	"os"
	"path"
	"strings"

	xhtml "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var debug *log.Logger
var info *log.Logger
var warning *log.Logger
var errorLog *log.Logger

func LogDefault() {
	debug = log.New(io.Discard, "DEBUG: ", log.LstdFlags|log.Lmsgprefix)
	info = log.New(io.Discard, "INFO: ", log.LstdFlags|log.Lmsgprefix)
	warning = log.New(os.Stderr, "WARNING: ", log.LstdFlags|log.Lmsgprefix)
	errorLog = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lmsgprefix)
}

func LogLevel(level string) error {
	LogDefault()
	switch level {
	case "debug":
		debug.SetOutput(os.Stderr)
		info.SetOutput(os.Stderr)
	case "info":
		info.SetOutput(os.Stderr)
	case "warning":
		// The default loglevel.
	case "error":
		warning.SetOutput(io.Discard)
	case "none":
		warning.SetOutput(io.Discard)
		errorLog.SetOutput(io.Discard)
	default:
		return errors.New("invalid log level")
	}
	return nil
}

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
			return err
		}
	} else {
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

func rewriteScript(node *xhtml.Node, file string) error {
	var text string

	if fh, err := os.Open(file); err == nil {
		if buf, err := io.ReadAll(fh); err == nil {
			text = string(buf)
		} else {
			return err
		}
	} else {
		return err
	}

	newChild := newTextNode(text)
	newAttrs := []xhtml.Attribute{}
	for _, at := range node.Attr {
		if at.Key != "src" {
			newAttrs = append(newAttrs, at)
		}
	}
	node.Attr = newAttrs
	node.AppendChild(newChild)

	return nil
}

// imageText returns a data URI for file based on the file extension. If the
// file is an svg image, it is inserted html escaped into the data URI,
// otherwise the file is base64url encoded.
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

// rewrteFaviconLink replaces the href attribute of a Link node with a data URI
// of file. It also modifies node's type attribute to match file's type.
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

// rewriteIMG replaces the src attribute of node with the data URI of file
// contents as output by imageText.
func rewriteIMG(node *xhtml.Node, file string) {
	for i := range node.Attr {
		if node.Attr[i].Key == "src" {
			node.Attr[i].Val = imageText(file)
		}
	}
}

// rewritePath returns old if old is not a key in rewriteRules, otherwise it
// returns rewriteRules[old].
func rewritePath(rewriteRules map[string]string, old string) string {
	if href, ok := rewriteRules[old]; ok {
		info.Printf("Rewriting %s -> %s.\n", old, href)
		return href
	} else {
		return old
	}
}

type elemType uint16

const (
	noneTag elemType = iota
	styleTag
	scriptTag
	faviconTag
	imgTag
)

// elemTypeString returns the string representation of each elemType except
// noneTag. Otherwise, it returns "".
func elemTypeString(elemtype elemType) string {
	switch elemtype {
	case styleTag:
		return "style"
	case scriptTag:
		return "script"
	case faviconTag:
		return "favicon"
	case imgTag:
		return "img"
	}
	return ""
}

// parseElement looks at node and if it is an internal resource, returns the
// corresponding elemType and the source path. It requires node to be a non-nil
// ElementType *xhtml.Node. It returns (noneTag, "") if the element is not an
// internal resource.
func parseElement(node *xhtml.Node) (elemtype elemType, source string) {
	if node.DataAtom == atom.Link {
		var rel string
		var href string
		for _, attr := range node.Attr {
			switch attr.Key {
			case "rel":
				rel = attr.Val
			case "href":
				href = attr.Val
			}
		}
		if rel == "stylesheet" {
			return styleTag, href
		} else if strings.Contains(rel, "icon") {
			return faviconTag, href
		}
	} else if node.DataAtom == atom.Img || node.DataAtom == atom.Script {
		src := ""
		for _, attr := range node.Attr {
			if attr.Key == "src" {
				src = attr.Val
			}
		}
		if node.DataAtom == atom.Img {
			return imgTag, src
		} else {
			return scriptTag, src
		}
	}
	return noneTag, ""
}

// crawler recursively searches the HTML document for internal resources. It
// returns false if node should be removed and true otherwise.
func RewriteCrawler(node *xhtml.Node, stem string, rewriteRules map[string]string) bool {
	// Skip checking non-element nodes.
	if node.Type == xhtml.ElementNode {
		if elemtype, source := parseElement(node); elemtype != noneTag {
			info.Printf("Found a %s element with source='%s'.\n", elemTypeString(elemtype), source)
			if source = rewritePath(rewriteRules, source); source == "" {
				return false
			}
			filepath := path.Join(stem, source)

			debug.Printf("Calculated pathname: '%s'.\n", filepath)

			switch elemtype {
			case styleTag:
				if err := rewriteCSSLink(node, filepath); err != nil {
					warning.Println(err)
				}
			case scriptTag:
				if err := rewriteScript(node, filepath); err != nil {
					warning.Println(err)
				}
			case faviconTag:
				rewriteFaviconLink(node, filepath)
			case imgTag:
				rewriteIMG(node, filepath)
			}

			return true
		}
	}

	for child := node.FirstChild; child != nil; {
		// Skip non-element child nodes and remove child if crawler returns false.
		if child.Type == xhtml.ElementNode && !RewriteCrawler(child, stem, rewriteRules) {
			info.Printf("Removing <%s> element.\n", child.Data)
			child = child.NextSibling
			if child != nil {
				node.RemoveChild(child.PrevSibling)
			} else {
				node.RemoveChild(node.LastChild)
			}
		} else {
			child = child.NextSibling
		}
	}
	return true
}
