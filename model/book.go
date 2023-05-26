package model

import "fmt"

type Book struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

func (b *Book) String() string {
	return fmt.Sprintf("Id: %d,name: %s,Author: %s", b.Id, b.Name, b.Author)
}
