package main

import (
	"github.com/ArtemT/books2vk"
	"github.com/davecgh/go-spew/spew"
	"github.com/namsral/flag"
)

func main() {
	var (
		path 	string
	)
	flag.StringVar(&path, "file", "", "Input XLSX file")
	flag.Parse()

	f := books2vk.WithFile(path)
	books := f.Read()

	spew.Dump(books)
}
