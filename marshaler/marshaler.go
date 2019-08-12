package marshaler

import (
	"bytes"
	"reflect"
	"strconv"
)

type Marshaler struct {
}

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

		if fv.Kind() == reflect.Struct || fv.Kind() == reflect.Ptr || fv.Kind() == reflect.Interface {
			d, err := m.marshal(fv)
			if err != nil {
				return nil, err
			}

			data = append(data, d...)
			continue
		}

		runes, err := extractRunes(fv)
		if err != nil {
			return nil, err
		}

		// check rune and strip
		ft := v.Type().Field(i)
		tag := ft.Tag.Get("fixed")
		limit, err := strconv.ParseInt(tag, 10, 64)
		if err != nil {
			//todo error value not valid
			return nil, err
		}

		padding := int(limit) - len(runes)
		if padding < 0 {
			runes = runes[0:limit]
		} else {
			paddingRunes := bytes.Runes(bytes.Repeat([]byte(" "), padding))
			runes = append(runes, paddingRunes...)
		}

		data = append(data, runes...)
	}

	return data, nil
}

func extractRunes(v reflect.Value) ([]rune, error) {
	switch v.Kind() {
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return []rune(strconv.Itoa(int(v.Int()))), nil
	case reflect.Float32:
		return []rune(strconv.FormatFloat(v.Float(), 'f', 2, 32)), nil
	case reflect.Float64:
		return []rune(strconv.FormatFloat(v.Float(), 'f', 2, 64)), nil
	case reflect.String:
		return []rune(v.String()), nil
	}

	return nil, nil
}
