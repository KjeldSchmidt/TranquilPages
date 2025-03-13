package models

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title   string `json:"title"`
	Author  string `json:"author"`
	Comment string `json:"comment"`
	Rating  int    `json:"rating"`
}

func CompareBooks(expected, actual *Book) bool {
	ignoreGORMFields := cmpopts.IgnoreFields(Book{}, "ID", "CreatedAt", "UpdatedAt", "DeletedAt")
	return cmp.Equal(expected, actual, ignoreGORMFields)
}
