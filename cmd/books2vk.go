package main

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/ArtemT/books2vk"
	"github.com/namsral/flag"
	"github.com/plandem/xlsx"
)

func main() {
	const (
		opCol = 19
		infoCol = 20
	)
	var (
		file string
	)
	flag.StringVar(&file, "file", "", "Input file")
	flag.Parse()

	xlsx, err := xlsx.Open(file)
	if err != nil {
		panic(err)
	}
	defer(xlsx).Close()

	sh := xlsx.Sheet(0)
	modified := false

	for rows := sh.Rows(); rows.HasNext(); {
		_, row := rows.Next()
		if op := row.Cell(opCol).String(); len(op) > 0 {
			b := books.Book{}
			ref := reflect.TypeOf(b)
			for i := 0; i < ref.NumField(); i++ {
				field := ref.Field(i)
				tag := field.Tag.Get("xcol")
				xcol, err := strconv.Atoi(tag)
				if err != nil {
					fmt.Println("cannot convert tag:", err)
					continue
				}
				switch field.Type.Name() {
				case "string":
					reflect.ValueOf(&b).Elem().FieldByName(field.Name).SetString(row.Cell(xcol).String())
				case "int":
					int, err := row.Cell(xcol).Int()
					if err != nil {
						fmt.Println("cannot convert to int:", err)
						continue
					}
					reflect.ValueOf(&b).Elem().FieldByName(field.Name).SetInt(int64(int))
				}
			}
			fmt.Printf("%+v\n", b)
			// row.Cell(opCol).Clear()
			// row.Cell(infoCol).SetDate(time.Now())
			// modified = true
		}
	}
	if modified {
		xlsx.Save()
	}
}
