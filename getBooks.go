package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func getBooks(c *gin.Context) {
	c.Header("Content_Type", "application/json")

	rows, err := db.Query("SELECT isbn, author, genre, publication_date FROM bookbase")
	if(err != nil) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []book

	for rows.Next() {
		var b book
		if err := rows.Scan(&b.ISBN, &b.Author, &b.Genre, &b.Publication_date); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		books = append(books, b)
	}

	err = rows.Err()
	if(err != nil) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, books)
}