package fixedwidth

import (
	"reflect"
	"strconv"
)

const tagName = "fixed"

type tag struct{}

// getLimitFixedTag get the tag `fixed` of a struct field then convert to integer
// if fixed tag is valid true will be returned; otherwise, false will be returned
func (tag) getLimitFixedTag(field reflect.StructField) (int, bool) {
	t, ok := field.Tag.Lookup(tagName)
	if !ok {
		return 0, true
	}

	l, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		ok = false
	}
	return int(l), ok
}
