package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	
)

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