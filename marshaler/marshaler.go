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

	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			r, err := m.marshal(v.Index(i))
			if err != nil {
				return nil, err
			}

			if i != v.Len()-1 {
				r = append(r, '\n')
			}
			data = append(data, r...)
		}
		return data, nil
	}

	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, nil
	}

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		tag := v.Type().Field(i).Tag.Get("fixed")
		limit, _ := strconv.ParseInt(tag, 10, 64)

		if fv.Kind() == reflect.Ptr || fv.Kind() == reflect.Interface {
			fv = fv.Elem()
		}

		if fv.Kind() == reflect.Struct {
			d, err := m.marshal(fv)
			if err != nil {
				return nil, err
			}

			d = truncateRunes(limit, d)
			data = append(data, d...)
			continue
		}

		runes := extractRunes(fv)
		runes = truncateRunes(limit, runes)
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
