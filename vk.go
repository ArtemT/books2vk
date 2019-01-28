package books2vk

import (
	"strconv"

	"github.com/davecgh/go-spew/spew"
)

type Service struct {
	Owner	string
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

func (vk Service) Send(in chan Book) chan Book {
	out := make(chan Book)
	go func() {
		defer close(out)
		for b := range in {
			spew.Dump("Send: " + strconv.Itoa(b.Row))
			spew.Dump(b)
			switch b.GetOp() {
			case 1: // Publish
				b.MktId = 123123
			case 2: // Unpublish
				b.MktId = 0
			}
			out <- b
		}
	}()
	return out
}


