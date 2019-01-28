package books2vk

import (
	"strconv"

	"github.com/davecgh/go-spew/spew"
)

type VK struct {
	Secret string
	// ...
}

func NewVK(conf string) VK {
	// o := vkapi.Options{}
	// api := vkapi.New(o)
	// ...
	vk := VK{
		// ...
	}
	return vk
}

func (vk VK) Send(in chan Book) chan Book {
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


