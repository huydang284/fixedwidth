package fixedwidth

const spaceByte = byte(' ')

// Marshal: see Marshal method of Marshaler
func Marshal(v interface{}) ([]byte, error) {
	return NewMarshaler().Marshal(v)
}

// Unmarshal: see Unmarshal method of Unmarshaler
func Unmarshal(data []byte, v interface{}) error {
	return NewUnmarshaler().Unmarshal(data, v)
}
