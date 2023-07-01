package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	xhtml "golang.org/x/net/html"
)

func main() {
	var help, version, listOnly, checkOnly bool
	var outFile, logLevel, rewrites, rewriteFile, wd string

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
	flagset.BoolVar(&help, "help", false, "Print this help text and exit.")
	flagset.BoolVar(&version, "version", false, "Print version and exit.")
	flagset.BoolVar(&listOnly, "list", false, "List mode: print internal resources and exit.")
	flagset.BoolVar(&checkOnly, "check", false, "Check run: print errors without rewriting and exit.")
	flagset.StringVar(&outFile, "o", "-", "Output file: - represents stdout. File is overwritten.")
	flagset.StringVar(&logLevel, "loglevel", "warning", "Log level: debug, info, warning, error, none.")
	flagset.StringVar(&rewrites, "r", "{}", "Rewrite rules: A JSON encoded map of rewrite rules.")
	flagset.StringVar(&rewriteFile, "f", "", "Rewrite rule file: A JSON file to read rules from.")
	flagset.StringVar(&wd, "wd", "", "Working directory. Defaults to index.html path.")

	if err := flagset.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// Setup loggers.
	if err := LogLevel(logLevel); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	names := make([]string, 0)
	flagset.Visit(func(flag *flag.Flag) {
		names = append(names, fmt.Sprintf("%s=%s", flag.Name, flag.Value.String()))
	})
	debug.Printf("Flags given: %s\n", strings.Join(names, ", "))

	if help {
		flagset.Usage()
		os.Exit(0)
	} else if version {
		fmt.Println("0.0.0")
		os.Exit(0)
	}

	if flagset.NArg() < 1 {
		flagset.Usage()
		os.Exit(2)
	}
	if wd == "" {
		wd = path.Dir(flagset.Arg(0))
	}

	file, err := os.Open(flagset.Arg(0))
	if err != nil {
		errorLog.Fatal(err)
	}
	doc, err := xhtml.Parse(file)
	if err != nil {
		errorLog.Fatal(err)
	}

	if listOnly {
		resources := RewriteCrawlList(doc)
		for _, resource := range resources {
			fmt.Print(resource)
		}
		os.Exit(0)
	}

	rewriteRules := make(map[string]string)
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

	// Get and rewrite HTML tree.
	RewriteCrawler(doc, wd, rewriteRules, checkOnly)
	if checkOnly {
		os.Exit(0)
	}

	// Output or write updated HTML.
	var writer io.Writer
	if outFile == "-" {
		writer = os.Stdout
	} else if writer, err = os.Create(outFile); err != nil {
		errorLog.Fatal(err)
	}

	if err := xhtml.Render(writer, doc); err != nil {
		errorLog.Fatal(err)
	}
	fmt.Fprintln(writer)
}
