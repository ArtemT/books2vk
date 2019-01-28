package books2vk

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

type Book struct {
	Author      string    `xcol:"1"`
	Title       string    `xcol:"2"`
	Description string    `xcol:"3"`
	Price       int       `xcol:"11"`
	Pic         fName     `xcol:"15"`
	MktId       int       `xcol:"19"`
	Op          operation `xcol:"20"`
	Row         int
}
type operation int
type fName string

// A bodge, keep it in sync with above
const (
	IdCol = 19
	OpCol = 20
)

func (b *Book) SetValues(f func(int) string) {
	var u Book
	ref := reflect.TypeOf(u)
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)
		tag := field.Tag.Get("xcol")
		if len(tag) == 0 {
			continue
		}
		xcol, err := strconv.Atoi(tag)
		if err != nil {
			fmt.Printf("cannot convert tag %s to column index: %v", tag, err)
			continue
		}
		fieldType := field.Type.Name()
		switch fieldType {
		case "string":
			val := f(xcol)
			reflect.ValueOf(b).Elem().FieldByName(field.Name).SetString(val)
		case "int":
			if s := f(xcol); len(s) > 0 {
				val, err := strconv.Atoi(s)
				if err != nil {
					fmt.Println("cannot convert value to int:", err)
					continue
				}
				reflect.ValueOf(b).Elem().FieldByName(field.Name).SetInt(int64(val))
			}
		case "fName":
			r, _ := regexp.Compile("[0-9]+\\.JPG")
			m := r.FindStringSubmatch(f(xcol))
			if len(m) > 0 {
				reflect.ValueOf(b).Elem().FieldByName(field.Name).SetString(m[0])
			}
		case "operation":
			b.setOp(f(xcol))
		default:
			fmt.Printf("type %s is not supported\n", fieldType)
		}
	}
}

func (b *Book) setPic(s string) {
	b.Pic = fName(s)
}

func (b Book) getPic() string {
	return string(b.Pic)
}

func (b *Book) setOp(s string) {
	if len(s) > 0 {
		i, err := strconv.Atoi(string(s[0]))
		if err != nil {
			fmt.Printf("not started with a number: %s, %v\n", s, err)
		}
		b.Op = operation(i)
	}
}

func (b Book) GetOp() int {
	return int(b.Op)
}
