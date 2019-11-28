package marshaler

import (
	"bytes"
	"reflect"
	"strconv"
)

type Marshaler struct {
	r []rune
}

func New() Marshaler {
	return Marshaler{}
}

func (m *Marshaler) Marshal(i interface{}) ([]rune, error) {
	m.reset()
	err := m.marshal(reflect.ValueOf(i))
	return m.r, err
}

func (m *Marshaler) reset() {
	m.r = m.r[:0]
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
				m.r = append(m.r, '\n')
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

		var runes []rune
		if fv.Kind() == reflect.Struct {
			startOffset := len(m.r) - 1
			err := m.marshal(fv)
			if err != nil {
				return err
			}

			endOffset := len(m.r) - 1
			if limitInt > 0 && endOffset-startOffset > limitInt {
				// truncate redundant runes
				m.r = m.r[:startOffset+limitInt+1]
			}
		} else {
			runes = extractRunes(fv)
			runes = truncateRunes(limitInt, runes)
			m.r = append(m.r, runes...)
		}
	}

	return nil
}

func truncateRunes(limit int, runes []rune) []rune {
	if limit == 0 {
		return runes
	}

	padding := limit - len(runes)
	if padding < 0 {
		runes = runes[0:limit]
	} else {
		paddingRunes := bytes.Runes(bytes.Repeat([]byte(" "), padding))
		runes = append(runes, paddingRunes...)
	}
	return runes
}

func extractRunes(v reflect.Value) []rune {
	switch v.Kind() {
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return []rune(strconv.Itoa(int(v.Int())))
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return []rune(strconv.FormatUint(v.Uint(), 10))
	case reflect.Float32:
		return []rune(strconv.FormatFloat(v.Float(), 'f', 2, 32))
	case reflect.Float64:
		return []rune(strconv.FormatFloat(v.Float(), 'f', 2, 64))
	case reflect.String:
		return []rune(v.String())
	}

	return nil
}
