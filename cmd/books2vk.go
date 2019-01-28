package main

import (
	. "github.com/ArtemT/books2vk"
	"github.com/namsral/flag"
)

func main() {
	var (
		path 	string
	)
	flag.StringVar(&path, "file", "", "Input XLSX file")
	flag.Parse()

	f := OpenFile(path)
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
