package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
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

func getBooks(c *gin.Context) {
	c.Header("Content_Type", "application/json")

	rows, err := db.Query("SELECT isbn, author, genre, publication_date FROM bookbase")
	if(err != nil) {
		log.Println(err)
	}
	defer rows.Close()

	var books []book

	for rows.Next() {
		var b book
		if err := rows.Scan(&b.ISBN, &b.Author, &b.Genre, &b.Publication_date); err != nil {
			log.Println(err)
		}
		books = append(books, b)
	}

	err = rows.Err()
	if(err != nil) {
		log.Println(err)
	}

	c.IndentedJSON(http.StatusOK, books)
}

func getBookByID(c *gin.Context) {
	isbn := c.Param("isbn")

	var b book
	err := db.QueryRow("SELECT isbn, author, genre, publication_date FROM bookbase WHERE isbn = $1", isbn).Scan(&b.ISBN, &b.Author, &b.Genre, &b.Publication_date)
	if(err != nil) {
		log.Println(err)
	}

	if (b.ISBN == "") {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "book not found"})
	} else {
		c.IndentedJSON(http.StatusOK, b)
	}
}

func createBook(c *gin.Context) {
	var newBook book
	err := c.BindJSON(&newBook)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("INSERT INTO bookbase (isbn, author, genre, publication_date) VALUES ($1, $2, $3, $4) RETURNING *")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(newBook.ISBN, newBook.Author, newBook.Genre, newBook.Publication_date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No rows affected"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Book created successfully"})
}