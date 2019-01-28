package main

import (
	. "github.com/ArtemT/books2vk"
)

func main() {
	f := OpenFile()
	defer func() {
		f.Save()
		f.Close()
	}()

	vk := NewService("")

	in := f.Proceed()
	out := vk.Send(in)
	done := f.Update(out)

	<-done
}
