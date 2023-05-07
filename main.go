package main

import (
	"library/config"
	"library/handlers"
	"log"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

var bh *handlers.BookHandler

func main() {
	log.Println("Statrting Library")

	bh = &handlers.BookHandler{DB: config.DBC.DB}
	router := gin.Default()
	router.GET("/books", bh.GetAllBooks)
	router.POST("/book", bh.AddBook)
	log.Println("Library Stated ")
	router.Run("localhost:8080")

}
