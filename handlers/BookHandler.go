package handlers

import (
	"fmt"
	"library/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookHandler struct {
	DB *gorm.DB
}

func (bh *BookHandler) GetAllBooks(ctx *gin.Context) {
	var books []model.Book
	if result := bh.DB.Find(&books); result.Error != nil {
		fmt.Printf("error is %v\n", result.Error)
		ctx.IndentedJSON(http.StatusInternalServerError, result.Error)
	} else {
		fmt.Printf("result is %v\n", &books)
		ctx.IndentedJSON(http.StatusOK, books)
	}
}

func (bh *BookHandler) AddBook(ctx *gin.Context) {
	book := &model.Book{}
	if err := ctx.BindJSON(book); err != nil {
		log.Println(err)
		return
	}

	if result := bh.DB.Create(book); result.Error != nil {
		log.Println(result.Error)
		ctx.IndentedJSON(http.StatusInternalServerError, result.Error)
	} else {
		fmt.Printf("result is %v\n", book)
		ctx.IndentedJSON(http.StatusOK, book)
	}
}
