package Core

import "reflect"

var ColorEnum = newcolorEnum()

func newcolorEnum() *color {
	return &color{
		Red:    1,
		Yellow: 2,
		Blue:   3,
		Black:  4,
		None:   5,
	}
}

type color struct {
	Red    int
	Yellow int
	Blue   int
	Black  int
	None   int
}

func GetEnumName(enum interface{}, value int) string {
	v := reflect.ValueOf(enum).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if int(v.Field(i).Int()) == value {
			return t.Field(i).Name
		}
	}
	return "Unknown"
}
