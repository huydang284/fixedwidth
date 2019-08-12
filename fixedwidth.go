package fixedwidth

import "github.com/huydang284/fixedwidth/marshaler"

func Marshal(v interface{}) ([]rune, error) {
	return marshaler.New().Marshal(v)
}
