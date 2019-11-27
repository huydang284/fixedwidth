package marshaler

import (
	"bytes"
	"reflect"
	"strconv"
)

type Marshaler struct{}

func New() Marshaler {
	return Marshaler{}
}

func (m Marshaler) Marshal(i interface{}) ([]rune, error) {
	return m.marshal(reflect.ValueOf(i))
}

func (m Marshaler) marshal(v reflect.Value) ([]rune, error) {
	var data []rune

	vKind := v.Kind()
	if vKind == reflect.Slice {
		vLen := v.Len()
		for i := 0; i < vLen; i++ {
			r, err := m.marshal(v.Index(i))
			if err != nil {
				return nil, err
			}

			if i != vLen-1 {
				r = append(r, '\n')
			}
			data = append(data, r...)
		}
		return data, nil
	}

	if vKind == reflect.Ptr || vKind == reflect.Interface {
		v = v.Elem()
	}

	if vKind != reflect.Struct {
		return nil, nil
	}

	vType := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		tag := vType.Field(i).Tag.Get("fixed")
		limit, _ := strconv.ParseInt(tag, 10, 64)

		if fv.Kind() == reflect.Ptr || fv.Kind() == reflect.Interface {
			fv = fv.Elem()
		}

		var runes []rune
		if fv.Kind() == reflect.Struct {
			d, err := m.marshal(fv)
			if err != nil {
				return nil, err
			}

			runes = truncateRunes(limit, d)
		} else {
			runes = extractRunes(fv)
			runes = truncateRunes(limit, runes)
		}

		data = append(data, runes...)
	}

	return data, nil
}

func truncateRunes(limit int64, runes []rune) []rune {
	if limit == 0 {
		return runes
	}

	padding := int(limit) - len(runes)
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
