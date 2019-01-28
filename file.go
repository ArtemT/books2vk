package books2vk

import (
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/plandem/xlsx"
)

type File struct {
	path     string
	modified bool
	doc      xlsx.Spreadsheet
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

func (f *File) Proceed() chan Book {
	sh := f.doc.Sheet(0)
	ch := make(chan Book)
	go func() {
		defer close(ch)
		for rows := sh.Rows(); rows.HasNext(); {
			i, row := rows.Next()
			// Don't care if no operation is required
			if len(row.Cell(OpCol).String()) == 0 {
				continue
			}
			b := Book{Row: i}
			b.SetValues(func(col int) string {
				return row.Cell(col).String()
			})
			spew.Dump("Proceed: " + strconv.Itoa(b.Row))
			ch <- b
		}
	}()
	return ch
}

func (f *File) Update(in chan Book) chan struct{} {
	sh := f.doc.Sheet(0)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for b := range in {
			spew.Dump("Update: " + strconv.Itoa(b.Row))

			// Save/remove market_item_id from VK
			mCell := sh.Cell(IdCol, b.Row)
			if b.MktId > 0 {
				mCell.SetInt(b.MktId)
				sh.Cell(OpCol, b.Row).Clear()
				f.modified = true
			} else if mCellVal, _ := mCell.Int(); mCellVal > 0 {
				mCell.Clear()
				sh.Cell(OpCol, b.Row).Clear()
				f.modified = true
			}
		}
	}()
	return done
}
