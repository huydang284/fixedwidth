package fixedwidth

import (
	"github.com/huydang284/fixedwidth/marshaler"
	"github.com/huydang284/fixedwidth/unmarshaler"
)

var m *marshaler.Marshaler

func Marshal(v interface{}) ([]byte, error) {
	if m == nil {
		newMarshaler := marshaler.New()
		m = &newMarshaler
	}
	return m.Marshal(v)
}

func Unmarshal(data []rune, v interface{}) error {
	return unmarshaler.New().Unmarshal(data, v)
}
