package books2vk

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Book struct {
	Author      string    `xcol:"1"`
	Title       string    `xcol:"2"`
	Description string    `xcol:"3"`
	Price       int       `xcol:"11"`
	Pic         numjpg    `xcol:"15"`
	MktId       int       `xcol:"19"`
	Op          operation `xcol:"20"`
	Row         int
}
type operation string
type numjpg string

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
		case "numjpg":
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
	b.Pic = numjpg(s)
}

func (b Book) getPic() string {
	return string(b.Pic)
}

func (b *Book) setOp(s string) {
	if len(s) > 0 {
		sl := strings.Split(s, "|")
		if len(sl) < 2 {
			fmt.Printf("No operation string found: %s.\n", s)
		}
		b.Op = operation(sl[1])
	}
}

func (b Book) GetOp() string {
	return string(b.Op)
}
