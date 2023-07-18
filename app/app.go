package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetPrefix("tf-server: ")

	log.Print("Creating server.")

	router := gin.Default()

	router.ForwardedByClientIP = true
	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("router.SetTrustedProxies: %v", err)
	}

	router.Static("/", "./www")

	if err := router.RunTLS("0.0.0.0:8080", "/etc/ssl/certs/cert.pem", "/etc/ssl/certs/key.pem"); err != nil {
		log.Fatal(err)
	}
}
