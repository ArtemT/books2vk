package book

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type Book struct {
	op          operation
	Author      string
	Title       string
	Description string
	Price       int
	Picture     numJPG
	MktId       int
	row         rowIdx
}
type operation string
type numJPG string
type rowIdx int

func (b *Book) SetOp(s string) {
	if len(s) > 0 {
		sl := strings.Split(s, "|")
		if len(sl) < 2 {
			fmt.Printf("No operation string found: %s.\n", s)
			return
		}
		b.op = operation(sl[1])
	} else {
		b.op = operation("")
	}
}
func (b Book) GetOp() string {
	return string(b.op)
}

func (b *Book) SetPic(s string) {
	r, _ := regexp.Compile("[0-9]+\\.JPG")
	m := r.FindStringSubmatch(s)
	if len(m) > 0 {
		b.Picture = numJPG(m[0])
	} else {
		fmt.Printf("Invalid filename: %s.\n", s)
	}
}
func (b Book) GetPic() string {
	return string(b.Picture)
}

func (b *Book) SetRow(i int) {
	b.row = rowIdx(i)
}
func (b Book) GetRow() int {
	return int(b.row)
}

type rowReader interface {
	Int(int) int
	String(int) string
	RowIdx() int
}

func (b *Book) GetValues(reader rowReader) bool {
	opCol := getCol("op")
	op := reader.String(opCol)
	// Break external loop early if no operation is set
	if len(op) == 0 {
		return false
	}
	b.SetOp(reader.String(opCol))
	b.SetRow(reader.RowIdx())

	ref := reflect.TypeOf(Book{})
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)
		if field.Name == "op" || field.Name == "row" { // TODO: Find a way to exclude unexported fields.
			continue
		}
		fieldType := field.Type.Name()
		col := getCol(field.Name)

		switch fieldType {
		case "int":
			reflect.ValueOf(b).Elem().FieldByName(field.Name).SetInt(int64(reader.Int(col)))
		case "string":
			reflect.ValueOf(b).Elem().FieldByName(field.Name).SetString(reader.String(col))
		case "numJPG":
			b.SetPic(reader.String(col))
		default:
			fmt.Printf("Field %s of type %s is not supported\n", field.Name, fieldType)
		}
	}
	return len(b.Author) > 0 && len(b.Title) > 0 && b.Price > 0 // TODO: Find out which fields are mandatory.
}

type rowWriter interface {
	Int(int, int)
	String(int, string)
}

func (b Book) SetValues(writer rowWriter) {
	writer.Int(getCol("MktId"), b.MktId)
	if len(b.op) == 0 {
		writer.String(getCol("op"), "")
	}
}

/** Config **/

type config struct {
	col map[string]int
	// TODO: Add templates here.
}

var conf *config

type configLoader interface {
	Load() map[string]interface{}
}

func ConfigInit(cl configLoader) {
	if conf != nil {
		return
	}
	c := config{col: make(map[string]int)}
	for name, v := range cl.Load() { // viper.GetStringMap("cols")
		col, ok := v.(int)
		if !ok {
			fmt.Printf("Cannot convert column %s index value %v to integer.\n", name, v)
			continue
		}
		c.col[name] = col
	}
	spew.Println("book.ConfigInit()")
	conf = &c
}

func getCol(fieldName string) int {
	col, ok := conf.col[strings.ToLower(fieldName)]
	if !ok {
		fmt.Printf("No column index is configured for field %s.\n", fieldName)
		return 0
	}
	return col
}
