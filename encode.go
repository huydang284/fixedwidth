package fixedwidth

import (
	"reflect"
	"strconv"
	"sync"
	"unicode/utf8"
)

type Marshaler struct {
	// mux is used to prevent other goroutines using the same Marshaler
	mux sync.Mutex

	// b is an underlying slice of bytes of a Marshaler.
	// After each marshal, b is reused via reset method.
	// By reusing b, we can minimize number of allocations
	b []byte
	tag
}

func NewMarshaler() *Marshaler {
	return &Marshaler{}
}

// Marshal returns the fixed-width encoding of v.
//
// v should be struct or a slice of struct.
//
// Each field in a struct need to be defined a `fixed` tag.
// The `fixed` tag indicates the maximum width of current field.
// Example:
// type foo struct {
//     bar string `fixed:"5"`
// }
// // value longer than `fixed` tag
// a := foo{bar: "longerthan5"}
// m.Marshal(a) // will returned []byte("longe") - remain characters are truncated.
// // value shorter than `fixed` tag
// b := foo{bar: "less"}
// m.Marshal(b) // will returned []byte("less ") - a padding space added.
//
// If b is slice of struct, Marshal will return multi lines seperated by new line character (\n).
func (m *Marshaler) Marshal(v interface{}) ([]byte, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.reset()
	err := m.marshal(reflect.ValueOf(v))
	return m.b, err
}

// reset the underlying slice, reuse memory allocation
func (m *Marshaler) reset() {
	m.b = m.b[:0]
}

func (m *Marshaler) marshal(v reflect.Value) error {
	vKind := v.Kind()
	if vKind == reflect.Slice {
		vLen := v.Len()
		for i := 0; i < vLen; i++ {
			err := m.marshal(v.Index(i))
			if err != nil {
				return err
			}

			if i != vLen-1 {
				m.b = append(m.b, '\n')
			}
		}
		return nil
	}

	if vKind == reflect.Ptr || vKind == reflect.Interface {
		v = v.Elem()
	}

	if vKind != reflect.Struct {
		return nil
	}

	vType := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		limit := m.getLimitFixedTag(vType.Field(i))

		if fv.Kind() == reflect.Ptr || fv.Kind() == reflect.Interface {
			fv = fv.Elem()
		}

		startOffset := len(m.b)
		if fv.Kind() == reflect.Struct {
			err := m.marshal(fv)
			if err != nil {
				return err
			}
		} else {
			m.appendExtractedScalarValue(fv)
		}

		if limit > 0 {
			m.truncateOrAddPadding(limit, startOffset)
		}
	}

	return nil
}

func (m *Marshaler) appendExtractedScalarValue(v reflect.Value) {
	switch v.Kind() {
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		m.b = strconv.AppendInt(m.b, v.Int(), 10)
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		m.b = strconv.AppendUint(m.b, v.Uint(), 10)
	case reflect.Float32:
		m.b = strconv.AppendFloat(m.b, v.Float(), 'f', 2, 32)
	case reflect.Float64:
		m.b = strconv.AppendFloat(m.b, v.Float(), 'f', 2, 64)
	case reflect.String:
		m.b = append(m.b, v.String()...)
	}

	return
}

func (m *Marshaler) truncateOrAddPadding(limit, lowerBound int) {
	if limit == 0 {
		return
	}

	b := m.b[lowerBound:]
	totalRunes := utf8.RuneCount(b)
	padding := limit - totalRunes
	if padding == 0 {
		return
	}

	if padding < 0 {
		// exclude redundant bytes
		m.b = m.b[:lowerBound+getFirstInvalidRune(limit, b)-1]
		return
	}

	// append additional spaces
	for i := 0; i < padding; i++ {
		m.b = append(m.b, spaceByte)
	}
	return
}

func getFirstInvalidRune(noRunes int, b []byte) int {
	i := 1
	for noRunes > 0 {
		_, s := utf8.DecodeRune(b)
		i += s
		noRunes--
	}
	return i
}
