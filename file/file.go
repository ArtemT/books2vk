package file

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ArtemT/books2vk/book"
	"github.com/davecgh/go-spew/spew"
	"github.com/plandem/xlsx"
	"github.com/spf13/viper"
)

type File struct {
	doc      xlsx.Spreadsheet
	modified bool
}

var file *File

func OpenFile() *File {
	if file == nil {
		file = &File{
			modified: false,
		}
	}
	path := viper.GetString("file")
	if len(path) == 0 {
		log.Fatal("No file specified. Usage: book2vk --file=example.xlsx")
	}
	file.Open(path)
	return file
}

func (f *File) Open(path string) {
	d, err := xlsx.Open(path)
	if err != nil {
		panic(err)
	}
	f.doc = *d
}

func (f File) Save() {
	if f.modified {
		err := f.doc.Save()
		if err != nil {
			panic(err)
		}
	}
}

func (f File) Close() {
	err := f.doc.Close()
	if err != nil {
		panic(err)
	}
}

type rowReader struct {
	row *xlsx.Row
	idx int
}

func (r rowReader) Int(i int) int {
	v := r.row.Cell(i).String()
	if len(v) == 0 {
		return 0
	}
	vi, err := strconv.Atoi(v)
	if err != nil {
		fmt.Printf("cannot read value as integer: %v\n", err)
	}
	return vi
}

func (r rowReader) String(i int) string {
	return r.row.Cell(i).Value()
}

func (r rowReader) RowNum() int {
	return r.idx
}

func (f *File) Read() chan book.Book {
	sh := f.doc.Sheet(0)
	ch := make(chan book.Book)
	go func() {
		defer close(ch)
		for rows := sh.Rows(); rows.HasNext(); {
			i, r := rows.Next()
			b := book.Book{}
			reader := rowReader{row: r, idx: i}
			if !b.GetValues(reader) {
				continue
			}
			spew.Dump("Read: " + strconv.Itoa(b.GetRow()))
			ch <- b
		}
	}()
	return ch
}

type rowWriter struct {
	row *xlsx.Row
}

func (r rowWriter) Int(i int, v int) {
	if v == 0 {
		r.row.Cell(i).Clear()
	} else {
		r.row.Cell(i).SetInt(v)
	}
}

func (r rowWriter) String(i int, v string) {
	if len(v) == 0 {
		r.row.Cell(i).Clear()
	} else {
		r.row.Cell(i).SetString(v)
	}
}

func (f *File) Update(in chan book.Book) chan struct{} {
	sh := f.doc.Sheet(0)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for b := range in {
			spew.Dump("Update: " + strconv.Itoa(b.GetRow()))
			writer := rowWriter{row: sh.Row(b.GetRow())}
			b.SetValues(writer)
			f.modified = true
		}
	}()
	return done
}
