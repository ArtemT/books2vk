package books2vk

import (
	"strconv"

	"github.com/ArtemT/books2vk/book"
	"github.com/davecgh/go-spew/spew"
)

const (
	ApiVer          = "5.92"
	BooksCategoryID = 901
)

type Service struct {
	Owner  string
	Secret string
	// ...
}

func NewService(conf string) Service {
	// o := vkapi.Options{}
	// api := vkapi.New(o)
	// ...
	vk := Service{
		// ...
	}
	return vk
}

func (vk Service) Send(in chan book.Book) chan book.Book {
	out := make(chan book.Book)
	go func() {
		defer close(out)
		for b := range in {
			spew.Dump("Send: " + strconv.Itoa(b.GetRow()))
			switch b.GetOp() {
			case "publish":
				b.MktId = 123123
			case "unpublish":
				b.MktId = 0
			}
			b.SetOp("")
			out <- b
		}
	}()
	return out
}
