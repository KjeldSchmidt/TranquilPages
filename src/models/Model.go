package models

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Author    string             `bson:"author" json:"author"`
	Comment   string             `bson:"comment" json:"comment"`
	Rating    int                `bson:"rating" json:"rating"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

func CompareBooks(expected, actual *Book) bool {
	ignoreFields := cmpopts.IgnoreFields(Book{}, "ID", "CreatedAt", "UpdatedAt")
	return cmp.Equal(expected, actual, ignoreFields)
}
