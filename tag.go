package fixedwidth

import (
	"reflect"
	"strconv"
)

const tagName = "fixed"

type tag struct{}

func (tag) getLimitFixedTag(field reflect.StructField) int {
	t := field.Tag.Get(tagName)
	l, _ := strconv.ParseInt(t, 10, 64)
	return int(l)
}
