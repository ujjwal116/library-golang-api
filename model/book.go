package model

type Book struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
}
