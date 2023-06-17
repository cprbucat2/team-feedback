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
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.Static("/", "./www")
	router.Run("0.0.0.0:8080")
}
