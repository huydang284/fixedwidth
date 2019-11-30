package fixedwidth

func Marshal(v interface{}) ([]byte, error) {
	return NewMarshaler().Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return NewUnmarshaler().Unmarshal(data, v)
}
