package handlers

import (
	"fmt"
	"library/model"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Database interface {
	Model(value interface{}) Database
	Count(count *int64) Database
	Order(id string) Database
	Create(book interface{}) Database
	Limit(n int) Database
	Where(query interface{}, args ...interface{}) Database
}

type BookRepository interface {
	GetBooksByPageAndPageSize(page int, pageSize int, books *[]model.Book) error
	GetBooksCount(count *int64) error
	AddBook(book *model.Book) error
}
type BookHandler struct {
	BookRepo BookRepository
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type PagingResult struct {
	Page       int           `json:"page"`
	PageSize   int           `json:"pageSize"`
	TotalItems int           `json:"totalItems"`
	TotalPages int           `json:"totalPages"`
	Data       *[]model.Book `json:"data"`
}

func (pr *PagingResult) String() string {
	return fmt.Sprintf("Page: %d, PageSize: %d, TotalItems: %d, TotalPages: %d, Books: %v",
		pr.Page, pr.PageSize, pr.TotalItems, pr.TotalPages, pr.Data)
}

func (bh *BookHandler) GetBooks(ctx *gin.Context) {

	if page, err := strconv.Atoi(ctx.DefaultQuery("page", "0")); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	} else if pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10")); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	} else {
		books := &[]model.Book{}
		count := new(int64)
		if err := bh.BookRepo.GetBooksCount(count); err != nil {
			response := ErrorResponse{
				Error: err.Error(),
			}
			ctx.JSON(http.StatusBadRequest, response)
		} else if pageNotAvailable(page, pageSize, *count) {
			respose := &ErrorResponse{
				Error: "Page not found",
			}
			fmt.Println(&respose)
			ctx.JSON(http.StatusOK, respose)
		} else if err = bh.BookRepo.GetBooksByPageAndPageSize(page, pageSize, books); err != nil {
			fmt.Printf("error is %v\n", err)
			response := ErrorResponse{
				Error: err.Error(),
			}
			ctx.JSON(http.StatusBadRequest, response)
		} else {
			totalpages := int(math.Ceil(float64(*count) / float64(pageSize)))
			pagingResult := PagingResult{page, pageSize, int(*count), totalpages, books}
			fmt.Println(&pagingResult)
			ctx.JSON(http.StatusOK, pagingResult)
		}
	}

}

func pageNotAvailable(page, pageSize int, count int64) bool {
	totalpages := int(math.Ceil(float64(count) / float64(pageSize)))
	return page > totalpages
}

func (bh *BookHandler) AddBook(ctx *gin.Context) {
	book := &model.Book{}
	if err := ctx.BindJSON(book); err != nil {
		log.Println(err)
		return
	}
	if err := bh.BookRepo.AddBook(book); err != nil {
		response := ErrorResponse{
			Error: err.Error(),
		}
		log.Println(response)
		ctx.JSON(http.StatusInternalServerError, response)
	} else {
		fmt.Printf("result is %v\n", *book)
		ctx.JSON(http.StatusOK, book)
	}
}
