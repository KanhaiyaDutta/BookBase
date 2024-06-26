package main

import(
	"net/http"
	"github.com/gin-gonic/gin"
)

func getBookByID(c *gin.Context) {
	isbn := c.Param("isbn")

	var b book
	err := db.QueryRow("SELECT isbn, author, genre, publication_date FROM bookbase WHERE isbn = $1", isbn).Scan(&b.ISBN, &b.Author, &b.Genre, &b.Publication_date)
	if(err != nil) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if (b.ISBN == "") {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "book not found"})
	} else {
		c.IndentedJSON(http.StatusOK, b)
	}
}