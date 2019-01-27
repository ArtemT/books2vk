package books2vk

import (
	"fmt"
	"reflect"
	"strconv"
)

type Book struct {
	Author      string    `xcol:"1"`
	Title       string    `xcol:"2"`
	Description string    `xcol:"3"`
	Price       int       `xcol:"11"`
	Op          operation `xcol:"19"`
	Status      string    `xcol:"20"`
	Row			int
}
// A bodge, keep it in sync with above
const OpCol = 19

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
		case "operation":
			b.setOp(f(xcol))
		default:
			fmt.Printf("type %s is not supported", fieldType)
		}
	}
}

type operation int

func (b *Book) setOp(s string) {
	if len(s) > 0 {
		i, err := strconv.Atoi(string(s[0]))
		if err != nil {
			fmt.Printf("not started with a number: %s, %v\n", s, err)
		}
		b.Op = operation(i)
	}
}
