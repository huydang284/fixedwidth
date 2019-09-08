package unmarshaler

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type Unmarshaler struct {
}

func New() Unmarshaler {
	return Unmarshaler{}
}

func (m Unmarshaler) Unmarshal(data []rune, model interface{}) error {
	modelType := reflect.TypeOf(model)

	if modelType.Kind() != reflect.Ptr {
		return errors.New("the model must be a pointer")
	}

	return m.unmarshal(data, reflect.ValueOf(model).Elem())
}

func (m Unmarshaler) unmarshal(data []rune, modelValue reflect.Value) error {
	modelType := modelValue.Type()

	if isBasicType(modelType.Kind()) {
		return m.unmarshalBasicType(data, modelValue)
	}

	switch modelType.Kind() {
	case reflect.Struct:
		return m.unmarshalStruct(data, modelValue)
	case reflect.Ptr:
		return m.unmarshalPointer(data, modelValue)
	case reflect.Slice:
		return m.unmarshalSlice(data, modelValue)
	case reflect.Interface:
		return m.unmarshalInterface(data, modelValue)
	}

	return nil
}

func (m Unmarshaler) unmarshalStruct(data []rune, modelValue reflect.Value) error {
	modelType := modelValue.Type()
	if modelType.Kind() != reflect.Struct {
		return errors.New("invalid type")
	}

	index := 0
	for i := 0; i < modelValue.NumField(); i++ {
		fieldType := modelType.Field(i)
		tag := fieldType.Tag.Get("fixed")
		l, _ := strconv.ParseInt(tag, 10, 64)
		limit := int(l)
		if limit == 0 {
			continue
		}

		pivot := index + limit
		if pivot > len(data) {
			pivot = len(data)
		}

		err := m.unmarshal(data[index:pivot], modelValue.Field(i))
		if err != nil {
			return err
		}

		// end of data
		if limit == len(data) {
			return nil
		}
		index += limit
	}

	return nil
}

func (m Unmarshaler) unmarshalBasicType(data []rune, modelValue reflect.Value) error {
	data = removePadding(data)
	modelType := modelValue.Type()

	switch modelType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(string(data), 10, 0)
		if err != nil {
			return err
		}
		modelValue.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(string(data), 10, 0)
		if err != nil {
			return err
		}
		modelValue.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return err
		}
		modelValue.SetFloat(f)
	case reflect.String:
		modelValue.SetString(string(data))
	}

	return nil
}

func (m Unmarshaler) unmarshalPointer(data []rune, modelValue reflect.Value) error {
	modelType := modelValue.Type()
	if modelType.Kind() != reflect.Ptr {
		return errors.New("invalid type")
	}

	newValue := reflect.New(modelValue.Type().Elem())
	err := m.unmarshal(data, newValue)
	if err != nil {
		return err
	}
	modelValue.Set(newValue)
	return nil
}

func (m Unmarshaler) unmarshalSlice(data []rune, modelValue reflect.Value) error {
	modelType := modelValue.Type()
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		newElem := reflect.New(modelType.Elem()).Elem()
		err := m.unmarshal([]rune(line), newElem)
		if err != nil {
			return err
		}
		modelValue.Set(reflect.Append(modelValue, newElem))
	}

	return nil
}

func (m Unmarshaler) unmarshalInterface(data []rune, modelValue reflect.Value) error {
	return m.unmarshal(data, modelValue.Elem())
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
