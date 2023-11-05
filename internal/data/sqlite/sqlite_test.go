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
	newLink := &linkzapp.Link{Name: "test", Url: "test.com", Labels: "testlabel", CreatedAt: 1234567890}
	db.Insert(newLink)
}
