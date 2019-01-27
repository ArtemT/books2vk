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

func (vk VK) Publish(in chan Book) chan Book {
	out := make(chan Book)
	go func() {
		defer close(out)
		for b := range in {
			spew.Dump("Publish: " + strconv.Itoa(b.Row))
			out <- b
		}
	}()
	return out
}


