package bookrepo

import (
	"library/model"

	"gorm.io/gorm"
)

type BookRepo struct {
	DB *gorm.DB
}

func (repo *BookRepo) GetBooksByPageAndPageSize(page int, pageSize int, books *[]model.Book) error {
	if res := repo.DB.Order("id").Limit(pageSize).Where("id >?", page*pageSize).Find(books); res.Error != nil {
		return res.Error
	} else {
		return nil
	}

}

func (repo *BookRepo) GetBooksCount(count *int64) error {
	q := repo.DB.Model(&model.Book{})
	res := q.Count(count)
	if res.Error != nil {
		return res.Error
	} else {
		return nil
	}

}

func (repo *BookRepo) AddBook(book *model.Book) error {
	if res := repo.DB.Create(book); res.Error != nil {
		return res.Error
	} else {
		return nil
	}
}
