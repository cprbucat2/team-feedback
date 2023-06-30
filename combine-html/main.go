package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	xhtml "golang.org/x/net/html"
)

func main() {
	// TODO: Process argument to write in place.
	// TODO: Process argument comtaining source rewrites -- FILE or text.
	// TODO: Add option to specify working directory other than index html stem.
	// TODO: Add -l, --list option to list external resources.
	// TODO: Add -s, --dry-run to check which external resources are missing.
	// TODO: Add log level option.
	var listOnly bool
	var checkOnly bool
	var outFile string
	var logLevel string
	var rewrites string
	var rewriteFile string

	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagset.Usage = func() {
		fmt.Fprintf(flagset.Output(), "Usage: %s [options] index.html\nOptions:\n", flagset.Name())
		flagset.PrintDefaults()
		fmt.Fprintln(flagset.Output(), `Rewrite rules:
  Rewrite rules are read from the command line specified with -r or from a file
with -f. Rules given with -f are read first and overridden by rules given with
-r. Rules are then applied in a single pass. Rules cannot rewrite other rules.
They should be provided as a JSON object mapping each URI/href to a local file
or empty string "".`)
	}
	flagset.BoolVar(&listOnly, "l", false, "List mode: print internal resources and exit.")
	flagset.BoolVar(&checkOnly, "c", false, "Check run: print errors without rewriting and exit.")
	flagset.StringVar(&outFile, "i", "-", "Output file: - represents stdout. File is overwritten.")
	flagset.StringVar(&logLevel, "loglevel", "warning", "Log level: debug, info, warning, error, none.")
	flagset.StringVar(&rewrites, "r", "{}", "Rewrite rules: A JSON encoded map of rewrite rules.")
	flagset.StringVar(&rewriteFile, "f", "", "Rewrite rule file: A JSON file to read rules from.")

	flagset.Parse(os.Args[1:])

	if flagset.NArg() < 1 {
		flagset.Usage()
		os.Exit(2)
	}
	wd := path.Dir(flagset.Arg(0))

	// Setup loggers.
	LogLevel(logLevel)

	file, err := os.Open(flagset.Arg(0))
	if err != nil {
		errorLog.Fatal(err)
	}
	doc, err := xhtml.Parse(file)
	if err != nil {
		errorLog.Fatal(err)
	}

	var rewriteRules map[string]string
	if rewriteFile != "" {
		if file, err := os.Open(rewriteFile); err != nil {
			errorLog.Fatal(err)
		} else {
			if buf, err := io.ReadAll(file); err != nil {
				errorLog.Fatal(err)
			} else {
				if err := json.Unmarshal(buf, &rewriteRules); err != nil {
					errorLog.Fatal(err)
				}
			}
		}
	}
	if rewrites != "" && rewrites != "{}" {
		var extraRewrites map[string]string
		if err := json.Unmarshal([]byte(rewrites), &extraRewrites); err != nil {
			errorLog.Fatal(err)
		}
		for k, v := range extraRewrites {
			rewriteRules[k] = v
		}
	}

	if listOnly {
		resources := RewriteCrawlList()
		for resource := range resources {
			fmt.Println(resource)
		}
		os.Exit(0)
	}

	// Get and rewrite HTML tree.
	RewriteCrawler(doc, wd, rewriteRules, checkOnly)
	if checkOnly {
		os.Exit(0)
	}

	// Output or write updated HTML.
	var buffer bytes.Buffer
	writer := io.Writer(&buffer)
	if err := xhtml.Render(writer, doc); err != nil {
		errorLog.Fatal(err)
	}
	fmt.Println(buffer.String())
}
