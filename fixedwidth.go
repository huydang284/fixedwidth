package fixedwidth

import (
	"github.com/huydang284/fixedwidth/marshaler"
	"github.com/huydang284/fixedwidth/unmarshaler"
)

func Marshal(v interface{}) ([]rune, error) {
	return marshaler.New().Marshal(v)
}

func Unmarshal(data []rune, v interface{}) error {
	return unmarshaler.New().Unmarshal(data, v)
}
