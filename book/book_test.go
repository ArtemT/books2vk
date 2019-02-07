package book

import (
	"reflect"
	"strconv"
	"testing"
)

type testConfLoader struct{}

func (testConfLoader) Load() map[string]interface{} {
	return map[string]interface{}{
		"op":          0,
		"author":      1,
		"title":       2,
		"description": 3,
		"price":       4,
		"picture":     5,
		"mktid":       6,
	}
}

func TestConfigInit(t *testing.T) {
	cl := testConfLoader{}
	ConfigInit(cl)
	for fieldName, i := range cl.Load() {
		ti := getCol(fieldName)
		if ti != i {
			t.Errorf("Field %s column index: actual %d, expected %d", fieldName, ti, i)
		}
	}
}

type TestRowReader struct {
	idx int
	row []string
}

func (tr TestRowReader) Int(i int) int {
	v, _ := strconv.Atoi(tr.row[i])
	return v
}
func (tr TestRowReader) String(i int) string {
	return tr.row[i]
}
func (tr TestRowReader) RowIdx() int {
	return tr.idx
}

type ReadTestResult struct {
	book    Book
	success bool
}

func GetReadTestResult(r TestRowReader) ReadTestResult {
	res := ReadTestResult{Book{}, false}
	res.success = res.book.GetValues(r)
	return res
}

func TestBook_GetValues(t *testing.T) {
	testValues := []TestRowReader{
		{0, []string{"Test |publish", "A", "T", "D", "100", ".123123.JPG.", ""}},
		{1, []string{"Test |unpublish", "A", "T", "D", "101", ".123123.JPG.", "123123"}},
		{2, []string{""}},
	}
	expected := []ReadTestResult{
		{
			Book{row: rowIdx(0), op: operation("publish"), Author: "A", Title: "T",
				Description: "D", Price: 100, Picture: "123123.JPG", MktId: 0},
			true,
		},
		{
			Book{row: rowIdx(1), op: operation("unpublish"), Author: "A", Title: "T",
				Description: "D", Price: 101, Picture: "123123.JPG", MktId: 123123},
			true,
		},
		{Book{}, false},
	}
	for i, r := range testValues {
		res := GetReadTestResult(r)
		if res != expected[i] {
			t.Errorf("\nActual:\n%v\nExpected:\n%v\n", res, expected[i])
		}
	}
}

type TestRowWriter struct {
	row map[int]string
}

func (tw *TestRowWriter) Int(col int, i int) {
	tw.row[col] = strconv.Itoa(i)
}
func (tw *TestRowWriter) String(col int, s string) {
	tw.row[col] = s
}

func TestBook_SetValues(t *testing.T) {
	testValues := []Book{
		{op: operation(""), Author: "A", Title: "T", Description: "D", Price: 100, Picture: "123123.JPG", MktId: 123123},
	}
	expected := []map[int]string{
		{0: "", 6: "123123"},
	}
	for i, r := range testValues {
		writer := &TestRowWriter{map[int]string{}}
		r.SetValues(writer)
		if !reflect.DeepEqual(writer.row, expected[i]) {
			t.Errorf("\nActual:\n%v\nExpected:\n%v\n", writer.row, expected[i])
		}
	}
}
