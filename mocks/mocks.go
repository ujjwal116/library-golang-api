package mocks

import (
	"library/model"

	"github.com/stretchr/testify/mock"
)

type BookRepoMock struct {
	mock.Mock
}

type MockResponseWriter struct {
	mock.Mock
}

func (mock *BookRepoMock) GetBooksByPageAndPageSize(page int, pageSize int, books *[]model.Book) error {
	args := mock.Called(page, pageSize, books)
	return args.Error(0)
}

func (mock *BookRepoMock) GetBooksCount(count *int64) error {
	args := mock.Called(count)
	return args.Error(0)
}

func (mock *BookRepoMock) AddBook(book *model.Book) error {
	args := mock.Called(book)
	return args.Error(0)
}
