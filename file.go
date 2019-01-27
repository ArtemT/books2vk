package books2vk

import (
	"github.com/plandem/xlsx"
)

type File struct {
	path     	string
	modified 	bool
	doc			xlsx.Spreadsheet
}

func OpenFile(p string) File {
	var f File
	f.path = p
	f.modified = false
	f.Open()
	return f
}

func (f *File) Open() {
	d, err := xlsx.Open(f.path)
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

func (f *File) Read() []Book {
	var books []Book
	sh := f.doc.Sheet(0)
	for rows := sh.Rows(); rows.HasNext(); {
		i, row := rows.Next()
		// Don't care if no operation is required
		if len(row.Cell(OpCol).String()) == 0 {
			break
		}
		b := Book{ Row: i }
		b.SetValues(func (col int) string {
			return row.Cell(col).String()
		})
		books = append(books, b)
	}
	return books
}
