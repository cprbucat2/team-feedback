package combine

import (
	"bytes"
	xhtml "golang.org/x/net/html"
	"regexp"
	"strings"
	"testing"
)

func TestRewriteCSSLink(t *testing.T) {
	style := `<html><head><link rel="stylesheet" href="styles.css"></head><body></body></html>`
	doc, err := xhtml.Parse(strings.NewReader(style))
	if err != nil {
		t.Fatal(err)
	}

	linkNode := doc.FirstChild.FirstChild.FirstChild

	if err := rewriteCSSLink(linkNode, "./rewrite_test-styles.css", false); err != nil {
		t.Error(err)
	}
	var buf bytes.Buffer
	if err := xhtml.Render(&buf, doc); err != nil {
		t.Error(err)
	}

	wantLink := regexp.MustCompile(`<link`)
	wantStyle := regexp.MustCompile(`<style>`)
	wantCSS := regexp.MustCompile(`.header { font-size: 2em; }`)
	if wantLink.Match(buf.Bytes()) || !wantStyle.Match(buf.Bytes()) || !wantCSS.Match(buf.Bytes()) {
		t.Errorf("Did not find match in %v.", buf.String())
	}
}
