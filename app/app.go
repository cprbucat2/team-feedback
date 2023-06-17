package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetPrefix("tf-server: ")

	log.Print("Connected to mysql.")

	router := gin.Default()
	router.Static("/", "./www")
	router.Run("0.0.0.0:8080")
}
