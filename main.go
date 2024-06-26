package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

type book struct {
	ISBN string `json:"isbn"`
	Author string `json:"author"`
	Genre string `json:"genre"`
	Publication_date string `json:"year"`
}

var db *sql.DB

func main() {
	godotenv.Load()
	var err error
	db, err = sql.Open("postgres", os.Getenv("AWS_ENDPOINT"))

	if(err != nil) {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n",err)
	}
	
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:isbn", getBookByID)
	router.POST("/books", createBook)

	router.Run("localhost:8080")
}