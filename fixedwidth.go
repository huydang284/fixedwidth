package fixedwidth

const spaceByte = byte(' ')

// see Marshal method of Marshaler
func Marshal(v interface{}) ([]byte, error) {
	return NewMarshaler().Marshal(v)
}

// see Unmarshal method of Unmarshaler
func Unmarshal(data []byte, v interface{}) error {
	return NewUnmarshaler().Unmarshal(data, v)
}
