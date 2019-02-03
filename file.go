package books2vk

import (
	"log"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/plandem/xlsx"
	"github.com/spf13/viper"
)

type File struct {
	doc      xlsx.Spreadsheet
	config   *viper.Viper
	modified bool
}

func OpenFile() File {
	f := File{
		config:   viper.Sub("xlsx"),
		modified: false,
	}
	path := f.config.GetString("file")
	if len(path) == 0 {
		log.Fatal("No file specified. Usage: book2vk --file=example.xlsx")
	}
	f.Open(path)
	return f
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

func (f *File) Proceed() chan Book {
	sh := f.doc.Sheet(0)
	ch := make(chan Book)
	opCol := f.config.GetInt("opcol")
	go func() {
		defer close(ch)
		for rows := sh.Rows(); rows.HasNext(); {
			i, row := rows.Next()
			// Don't care if no operation is required
			if len(row.Cell(opCol).String()) == 0 {
				continue
			}
			b := Book{Row: i}
			// @todo: It fails sometimes. Make a test.
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
	opCol := f.config.GetInt("opcol")
	idCol := f.config.GetInt("idcol")
	go func() {
		defer close(done)
		for b := range in {
			spew.Dump("Update: " + strconv.Itoa(b.Row))

			// Save/remove market_item_id
			// @todo: Move it to callback in vk.go
			mCell := sh.Cell(idCol, b.Row)
			if b.MktId > 0 {
				mCell.SetInt(b.MktId)
				sh.Cell(opCol, b.Row).Clear()
				f.modified = true
			} else if mCellVal, _ := mCell.Int(); mCellVal > 0 {
				mCell.Clear()
				sh.Cell(opCol, b.Row).Clear()
				f.modified = true
			}
		}
	}()
	return done
}
