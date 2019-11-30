package fixedwidth

import (
	"bytes"
	"errors"
	"reflect"
	"strconv"
	"unicode/utf8"
)

type Unmarshaler struct {
	tag
}

func NewUnmarshaler() Unmarshaler {
	return Unmarshaler{}
}

func (m Unmarshaler) Unmarshal(data []byte, model interface{}) error {
	modelType := reflect.TypeOf(model)

	if modelType.Kind() != reflect.Ptr {
		return errors.New("the model must be a pointer")
	}

	_, err := m.unmarshal(data, reflect.ValueOf(model).Elem(), reflect.TypeOf(model).Elem())
	return err
}

func (m Unmarshaler) unmarshal(data []byte, modelValue reflect.Value, modelType reflect.Type) (int, error) {
	if isBasicType(modelType.Kind()) {
		return m.unmarshalBasicType(data, modelValue)
	}

	switch modelType.Kind() {
	case reflect.Struct:
		return m.unmarshalStruct(data, modelValue)
	case reflect.Ptr:
		return m.unmarshalPointer(data, modelValue, modelType)
	case reflect.Slice:
		return m.unmarshalSlice(data, modelValue)
	case reflect.Interface:
		return m.unmarshalInterface(data, modelValue)
	}

	return 0, nil
}

func (m Unmarshaler) unmarshalStruct(data []byte, structValue reflect.Value) (int, error) {
	structType := structValue.Type()
	if structType.Kind() != reflect.Struct {
		return 0, errors.New("input value is not a struct")
	}

	index := 0
	dataLen := len(data)
	for i := 0; i < structValue.NumField(); i++ {
		if index >= dataLen {
			break
		}

		structField := structType.Field(i)
		fieldValue := structValue.Field(i)
		limit := m.getLimitFixedTag(structField)

		if limit == 0 && !isStructOrStructPointer(structField.Type) {
			continue
		}

		var uLen int
		var err error
		if limit == 0 && isStructOrStructPointer(structField.Type) {
			uLen, err = m.unmarshal(data[index:], fieldValue, structField.Type)
		} else {
			upperBound := getUpperBound(index, limit, data)
			uLen, err = m.unmarshal(data[index:upperBound], fieldValue, structField.Type)
		}

		if err != nil {
			return 0, err
		}

		index += uLen
	}

	return index, nil
}

func (m Unmarshaler) unmarshalBasicType(data []byte, modelValue reflect.Value) (int, error) {
	l := len(data)
	data = removePadding(data)
	modelType := modelValue.Type()

	switch modelType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(string(data), 10, 0)
		if err != nil {
			return 0, err
		}
		modelValue.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(string(data), 10, 0)
		if err != nil {
			return 0, err
		}
		modelValue.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return 0, err
		}
		modelValue.SetFloat(f)
	case reflect.String:
		modelValue.SetString(string(data))
	}

	return l, nil
}

func (m Unmarshaler) unmarshalPointer(data []byte, modelValue reflect.Value, modelType reflect.Type) (int, error) {
	if modelType.Kind() != reflect.Ptr {
		return 0, errors.New("invalid type")
	}

	newType := modelType.Elem()
	newValue := reflect.New(newType)
	l, err := m.unmarshal(data, newValue.Elem(), newType)
	if err != nil {
		return 0, err
	}
	modelValue.Set(newValue)
	return l, nil
}

func (m Unmarshaler) unmarshalSlice(data []byte, modelValue reflect.Value) (int, error) {
	modelType := modelValue.Type()
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		newElem := reflect.New(modelType.Elem()).Elem()
		_, err := m.unmarshal(line, newElem, modelType.Elem())
		if err != nil {
			return 0, err
		}
		modelValue.Set(reflect.Append(modelValue, newElem))
	}

	return len(data), nil
}

func (m Unmarshaler) unmarshalInterface(data []byte, modelValue reflect.Value) (int, error) {
	var tempString string
	newType := reflect.TypeOf(tempString)
	newValue := reflect.New(newType)
	l, err := m.unmarshal(data, newValue.Elem(), newType)
	if err != nil {
		return 0, err
	}
	modelValue.Set(newValue.Elem())
	return l, nil
}

func removePadding(data []byte) []byte {
	return bytes.Trim(data, " ")
}

func isBasicType(p reflect.Kind) bool {
	basicTypes := []reflect.Kind{
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String,
	}

	for _, t := range basicTypes {
		if t == p {
			return true
		}
	}

	return false
}

func isStructOrStructPointer(t reflect.Type) bool {
	if t.Kind() == reflect.Struct {
		return true
	}

	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		return true
	}

	return false
}

func getUpperBound(lowerBound, limit int, data []byte) int {
	diff := 0
	lenData := len(data)
	data = data[lowerBound:]
	for limit > 0 {
		_, s := utf8.DecodeRune(data)
		diff += s
		limit--
		data = data[s:]
	}
	upperBound := lowerBound + diff
	if upperBound > lenData {
		upperBound = lenData
	}

	return upperBound
}
