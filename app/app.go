package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"

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
	router.POST("/api/submit", postUserSubmission)

	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatal(err)
	}
}

type score struct {
	Participation float64 `json:"participation"`
	Collaboration float64 `json:"collaboration"`
	Contribution  float64 `json:"contribution"`
	Attitude      float64 `json:"attitude"`
	Goals         float64 `json:"goals"`
}

type groupData struct {
	Name    string    `json:"name"`
	Scores  []float64 `json:"scores"`
	Comment string    `json:"comment"`
}
type userSubmission struct {
	StudentName  string      `json:"studentName"`
	StudentGroup string      `json:"studentGroup"`
	GroupSize    int64       `json:"groupSize"`
	GroupData    []groupData `json:"groupsData"`
}

func postUserSubmission(c *gin.Context) {
	submission := userSubmission{}
	c.BindJSON(&submission)
	log.Println("postUserSubmission: got ", submission)
	subUUID := uuid.NewString()
	log.Println(subUUID)
	res := struct {
		Uuid string `json:"uuid"`
	}{subUUID}
	c.JSON(http.StatusOK, res)
}
