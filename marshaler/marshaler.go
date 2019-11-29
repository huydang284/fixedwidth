package marshaler

import (
	"bytes"
	"reflect"
	"strconv"
	"unicode/utf8"
)

type Marshaler struct {
	b []byte
}

func New() Marshaler {
	return Marshaler{}
}

func (m *Marshaler) Marshal(i interface{}) ([]byte, error) {
	m.reset()
	err := m.marshal(reflect.ValueOf(i))
	return m.b, err
}

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

		var limit int64
		tag := vType.Field(i).Tag.Get("fixed")
		if tag != "" {
			limit, _ = strconv.ParseInt(tag, 10, 64)
		}

		limitInt := int(limit)

		if fv.Kind() == reflect.Ptr || fv.Kind() == reflect.Interface {
			fv = fv.Elem()
		}

		var data []byte

		startOffset := len(m.b)
		if fv.Kind() == reflect.Struct {
			err := m.marshal(fv)
			if err != nil {
				return err
			}
		} else {
			data = extractRunes(fv)
			m.b = append(m.b, data...)
		}

		if limitInt > 0 {
			// truncate redundant runes
			m.truncate(limitInt, startOffset)
		}
	}

	return nil
}

func (m *Marshaler) truncate(limit, start int) {
	if limit == 0 {
		return
	}

	b := m.b[start:]
	totalRunes := utf8.RuneCount(b)
	padding := limit - totalRunes
	if padding == 0 {
		return
	}

	if padding < 0 {
		// exclude redundant bytes
		m.b = m.b[:start+getFirstInvalidRune(limit, b)-1]
		return
	}

	paddingBytes := bytes.Repeat([]byte(" "), padding)
	m.b = append(m.b, paddingBytes...)
	return
}

func getFirstInvalidRune(limit int, b []byte) int {
	i := 1
	for limit > 0 {
		_, s := utf8.DecodeRune(b)
		i += s
		limit--
	}
	return i
}

func extractRunes(v reflect.Value) []byte {
	switch v.Kind() {
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return []byte(strconv.Itoa(int(v.Int())))
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return []byte(strconv.FormatUint(v.Uint(), 10))
	case reflect.Float32:
		return []byte(strconv.FormatFloat(v.Float(), 'f', 2, 32))
	case reflect.Float64:
		return []byte(strconv.FormatFloat(v.Float(), 'f', 2, 64))
	case reflect.String:
		return []byte(v.String())
	}

	return nil
}
