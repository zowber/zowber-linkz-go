package sqlite

import (
	"testing"

	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

var db = SQLiteClient{}

func TestGet(t *testing.T) {
	//fmt.Println(db.All())
}

func TestInsert(t *testing.T) {
	newLink := &linkzapp.Link{}
	db.Insert(newLink)
}
