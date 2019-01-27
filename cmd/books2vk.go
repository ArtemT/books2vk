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

	vk := NewVK("")

	in := f.Proceed()
	out := vk.Publish(in)
	done := f.Update(out)

	<-done
}
