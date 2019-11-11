package unmarshaler

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type Unmarshaler struct{}

func New() Unmarshaler {
	return Unmarshaler{}
}

func (m Unmarshaler) Unmarshal(data []rune, model interface{}) error {
	modelType := reflect.TypeOf(model)

	if modelType.Kind() != reflect.Ptr {
		return errors.New("the model must be a pointer")
	}

	_, err := m.unmarshal(data, reflect.ValueOf(model).Elem(), reflect.TypeOf(model).Elem())
	return err
}

func (m Unmarshaler) unmarshal(data []rune, modelValue reflect.Value, modelType reflect.Type) (int, error) {
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

func (m Unmarshaler) unmarshalStruct(data []rune, modelValue reflect.Value) (int, error) {
	modelType := modelValue.Type()
	if modelType.Kind() != reflect.Struct {
		return 0, errors.New("invalid type")
	}

	index := 0
	for i := 0; i < modelValue.NumField(); i++ {
		if index >= len(data) {
			break
		}

		fieldType := modelType.Field(i)
		fieldValue := modelValue.Field(i)
		tag := fieldType.Tag.Get("fixed")
		l, _ := strconv.ParseInt(tag, 10, 64)
		limit := int(l)
		if limit == 0 && !isStructOrStructPointer(fieldType.Type) {
			continue
		}

		var err error
		var uLen int
		if limit == 0 && isStructOrStructPointer(fieldType.Type) {
			uLen, err = m.unmarshal(data[index:], fieldValue, fieldType.Type)
		} else {
			pivot := index + limit
			if pivot > len(data) {
				pivot = len(data)
			}

			uLen, err = m.unmarshal(data[index:pivot], fieldValue, fieldType.Type)
		}

		if err != nil {
			return 0, err
		}

		index += uLen
	}

	return index, nil
}

func (m Unmarshaler) unmarshalBasicType(data []rune, modelValue reflect.Value) (int, error) {
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

func (m Unmarshaler) unmarshalPointer(data []rune, modelValue reflect.Value, modelType reflect.Type) (int, error) {
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

func (m Unmarshaler) unmarshalSlice(data []rune, modelValue reflect.Value) (int, error) {
	modelType := modelValue.Type()
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		newElem := reflect.New(modelType.Elem()).Elem()
		_, err := m.unmarshal([]rune(line), newElem, modelType.Elem())
		if err != nil {
			return 0, err
		}
		modelValue.Set(reflect.Append(modelValue, newElem))
	}

	return len(data), nil
}

func (m Unmarshaler) unmarshalInterface(data []rune, modelValue reflect.Value) (int, error) {
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

// todo check this performance
func removePadding(data []rune) []rune {
	return []rune(strings.Trim(string(data), " "))
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
