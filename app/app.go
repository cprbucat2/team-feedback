package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func redirect(w http.ResponseWriter, req *http.Request) {
	log.Print(req.Host)
	req.URL.Host = req.Host[:len(req.Host)-3] + "443"
	req.URL.Scheme = "https"
	log.Print(req.URL.String())
	http.Redirect(w, req,
		req.URL.String(),
		http.StatusMovedPermanently)
}

func main() {
	// redirect every http request to https
	go func() {
		err := http.ListenAndServe(":8080", http.HandlerFunc(redirect))
		log.Fatal(err)
	}()

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

	if err := router.RunTLS("0.0.0.0:8443", "/etc/ssl/certs/cert.pem", "/etc/ssl/certs/key.pem"); err != nil {
		log.Fatal(err)
	}
}

type entry struct {
	Name    string    `json:"name"`
	Scores  []float64 `json:"scores"`
	Comment string    `json:"comment"`
}

type submission struct {
	Author      string  `json:"author"`
	Entries     []entry `json:"entries"`
	Improvement string  `json:"improvement"`
}

func postUserSubmission(c *gin.Context) {
	s := submission{}
	if err := c.BindJSON(&s); err != nil {
		log.Println("postUserSubmission:", err)
		c.Status(http.StatusBadRequest)
		return
	}
	log.Println("postUserSubmission: got", s)

	// TODO: Query teams table for team member IDs.

	var sub_id int64
	// TODO: Insert Author ID instead of 1!!!
	if result, err := db.Exec("insert into submissions (author, improvement) values (?, ?)",
		1, s.Improvement); err != nil {
		log.Println("postUserSubmission:", err)
		c.Status(http.StatusInternalServerError)
		return
	} else {
		sub_id, err = result.LastInsertId()
		if err != nil {
			log.Println("postUserSubmission:", err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	log.Println("postUserSubmission: inserted", sub_id)

	for _, e := range s.Entries {
		// TODO: Insert user ID of the member instead of 1!!!
		if _, err := db.Exec(`insert into entries (submission_id, member,
			Participation, Collaboration, Contribution, Attitude, Goals, Comment)
			values (?, ?, ?, ?, ?, ?, ?, ?)`, sub_id, 1, e.Scores[0], e.Scores[1],
			e.Scores[2], e.Scores[3], e.Scores[4], e.Comment); err != nil {
			log.Println("postUserSubmission:", err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": sub_id,
	})
}
