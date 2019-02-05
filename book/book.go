package book

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

type Book struct {
	Operation   operation
	Author      string
	Title       string
	Description string
	Price       int
	Picture     numJPG
	MktId       int
	row         rowNum
}
type operation string
type numJPG string
type rowNum int

func (b *Book) SetOp(s string) {
	if len(s) > 0 {
		sl := strings.Split(s, "|")
		if len(sl) < 2 {
			fmt.Printf("No operation string found: %s.\n", s)
			return
		}
		b.Operation = operation(sl[1])
	} else {
		b.Operation = ""
	}
}

func (b Book) GetOp() string {
	return string(b.Operation)
}

func (b *Book) SetPic(s string) {
	b.Picture = numJPG(s)
}

func (b Book) GetPic() string {
	return string(b.Picture)
}

func (b Book) SetRow(i int) {
	b.row = rowNum(i)
}

func (b Book) GetRow() int {
	return int(b.row)
}

type rowReader interface {
	Int(int) int
	String(int) string
	RowNum() int
}

func (b *Book) GetValues(reader rowReader) bool {
	ref := reflect.TypeOf(Book{})
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)
		if field.Name == "row" { // TODO: Exclude unexported fields.
			continue
		}
		col := getCol(field.Name)
		fieldType := field.Type.Name()
		spew.Println(field.Name)
		switch fieldType {
		case "operation":
			op := reader.String(col)
			// Break external loop early if no operation is set
			if len(op) == 0 {
				return false
			}
			b.SetOp(reader.String(col))
		case "int":
			reflect.ValueOf(b).Elem().FieldByName(field.Name).SetInt(int64(reader.Int(col)))
		case "string":
			reflect.ValueOf(b).Elem().FieldByName(field.Name).SetString(reader.String(col))
		case "numJPG":
			r, _ := regexp.Compile("[0-9]+\\.JPG")
			m := r.FindStringSubmatch(reader.String(col))
			if len(m) > 0 {
				reflect.ValueOf(b).Elem().FieldByName(field.Name).SetString(m[0])
			}
		case "rowNum":
			b.SetRow(reader.RowNum())
		default:
			fmt.Printf("type %s is not supported\n", fieldType)
		}
	}
	return len(b.Author) > 0 && len(b.Title) > 0 && b.Price > 0
}

type rowWriter interface {
	Int(int, int)
	String(int, string)
}

func (b Book) SetValues(writer rowWriter) {
	writer.Int(getCol("MktId"), b.MktId)
	if len(b.Operation) == 0 {
		writer.String(getCol("Operation"), "")
	}
}

/** Config **/

type config struct {
	col map[string]int
	// TODO: Add templates here.
}

var conf *config

func newConfig() *config {
	c := config{col: make(map[string]int)}
	for name, v := range viper.GetStringMap("cols") {
		col, ok := v.(int)
		if !ok {
			fmt.Printf("Cannot convert column %s index value %v to integer.\n", name, v)
			continue
		}
		c.col[name] = col
	}
	spew.Println("book.newConfig")
	return &c
}

func initConfig() {
	if conf == nil {
		conf = newConfig()
	}
}

func getCol(fieldName string) int {
	initConfig()
	col, ok := conf.col[strings.ToLower(fieldName)]
	if !ok || col == 0 {
		fmt.Printf("No column index is configured for field %s.\n", fieldName)
		return 0
	}
	return col
}
