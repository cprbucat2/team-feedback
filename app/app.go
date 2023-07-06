package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/cprbucat2/team-feedback/app/combine"
	"github.com/gin-gonic/gin"
	xhtml "golang.org/x/net/html"
)

func main() {
	log.SetPrefix("tf-server: ")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s <command> [arguments]\n\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), `Commands:

	serve			Run server.
	generate	Generate static HTML.`)
		flag.PrintDefaults()
	}

	if len(os.Args) < 2 {
		doServer()
	} else {
		switch os.Args[1] {
		case "serve":
			doServer()
		case "generate":
			doGenerate()
		default:
			flag.Usage()
			os.Exit(2)
		}
	}
}

func doServer() {
	log.Print("Creating server.")

	router := gin.Default()

	router.ForwardedByClientIP = true
	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("router.SetTrustedProxies: %v", err)
	}

	router.Static("/", "./www")

	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatal(err)
	}
}

func doGenerate() {
	wd := path.Join(path.Dir(os.Args[0]), "www")
	file, err := os.Open(path.Join(wd, "index.html"))
	if err != nil {
		log.Fatal(err)
	}
	doc, err := xhtml.Parse(file)
	if err != nil {
		log.Fatal(err)
	}
	rewrites := make(map[string]string)

	combine.LogDefault()
	combine.RewriteCrawler(doc, wd, rewrites, false)
	if err := xhtml.Render(os.Stdout, doc); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}
