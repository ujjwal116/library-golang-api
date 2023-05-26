package main

import (
	"library/bookrepo"
	"library/config"
	"library/handlers"
	"log"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

var bh *handlers.BookHandler

func main() {
	log.Println("Statrting Library")
	bookRepo := &bookrepo.BookRepo{DB: config.DBC.DB}
	bh = &handlers.BookHandler{BookRepo: bookRepo}
	router := gin.Default()
	router.GET("/books", bh.GetBooks)
	router.POST("/book", bh.AddBook)
	log.Println("Library Stated ")
	router.Run("localhost:8080")
}
