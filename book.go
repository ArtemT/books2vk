package books2vk

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"github.com/oklog/ulid"
)

type Book struct {
	Id          ulid.ULID `xcol:"21"`
	Author      string    `xcol:"1"`
	Title       string    `xcol:"2"`
	Description string    `xcol:"3"`
	Price       int       `xcol:"11"`
	Operation   op        `xcol:"19"`
	Status      string    `xcol:"20"`
}

func (b *Book) SetValues(f func(int) string) {
	var u Book
	ref := reflect.TypeOf(u)
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)
		xcol, err := strconv.Atoi(field.Tag.Get("xcol"))
		if err != nil {
			fmt.Println("cannot convert tag to column index:", err)
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
		case "op":
			b.Operation.setId(f(xcol))
		case "ULID":
			b.setId(f(xcol))
		default:
			fmt.Printf("type %s is not supported", fieldType)
		}
	}
}

func (b *Book) setId(s string) {
	if len(s) > 0 {
		if u, err := ulid.ParseStrict(s); err == nil {
			b.Id = u
			return
		} else {
			fmt.Println("wrong ULID:", err)
		}
	}
	t := time.Unix(1000000, 0)
	e := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	b.Id = ulid.MustNew(ulid.Timestamp(t), e)
}

func (b Book) getId() string {
	return b.Id.String()
}

type op struct {
	OpId int
}

func (o *op) setId(s string) {
	if len(s) > 0 {
		i, err := strconv.Atoi(string(s[0]))
		if err != nil {
			fmt.Printf("not started with a number: %s, %v\n", s, err)
		}
		o.OpId = i
	}
}
