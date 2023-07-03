package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	log.SetPrefix("tf-server: ")

	dbConfig := mysql.Config{
		User:   os.Getenv("MYSQL_USER"),
		Passwd: os.Getenv("MYSQL_PASSWORD"),
		Net:    "tcp",
		Addr:   "mysql:3306",
		DBName: os.Getenv("MYSQL_DB"),
	}

	var dbOpenErr error
	if db, dbOpenErr = sql.Open("mysql", dbConfig.FormatDSN()); dbOpenErr != nil {
		log.Fatal(dbOpenErr)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Print("Connected to MySQL.")
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
	if err := c.BindJSON(&submission); err != nil {
		log.Println("ERROR: ", err)
	}
	log.Println("postUserSubmission: got ", submission)
	subUUID := uuid.NewString()
	log.Println(subUUID)

	// TODO: Query teams table for team member IDs.

	submissionIDs := make([]string, 0)
	for i, memberData := range submission.GroupData {
		result, err := db.Exec("insert into submissions (Member, Participation, Collaboration, Contribution, Attitude, Goals, Comment) values (?, ?, ?, ?, ?, ?, ?)",
			i, memberData.Scores[0], memberData.Scores[1], memberData.Scores[2], memberData.Scores[3], memberData.Scores[4], memberData.Comment)
		if err != nil {
			log.Printf("postUserSubmission: %v\n", err)
		}
		if id, err := result.LastInsertId(); err != nil {
			log.Printf("postUserSubmission: %v\n", err)
		} else {
			submissionIDs = append(submissionIDs, strconv.FormatInt(id, 10))
		}
	}

	submissionIDsString := strings.Join(submissionIDs, ",")

	result, err := db.Exec("insert into membersubmissions (UUID, Author, Submissions, Improvement) values (uuid_to_bin(?), 1, ?, ?)",
		subUUID, submissionIDsString, submission.GroupData[0].Comment)
	if err != nil {
		log.Println("postUserSubmission: ", err)
	}
	subID, err := result.LastInsertId()
	if err != nil {
		log.Println("postUserSubmission: ", err)
	}

	res := struct {
		Id   int64  `json:"id"`
		Uuid string `json:"uuid"`
	}{subID, subUUID}

	c.JSON(http.StatusOK, res)
}
