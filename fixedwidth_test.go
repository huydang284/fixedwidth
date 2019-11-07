package fixedwidth

func float64p(v float64) *float64 { return &v }
func stringp(v string) *string    { return &v }
func intp(v int) *int             { return &v }
func int8p(v int8) *int8          { return &v }
func int16p(v int16) *int16       { return &v }
func int32p(v int32) *int32       { return &v }
func int64p(v int64) *int64       { return &v }
func uintp(v uint) *uint          { return &v }
func uint8p(v uint8) *uint8       { return &v }
func uint16p(v uint16) *uint16    { return &v }
func uint32p(v uint32) *uint32    { return &v }
func uint64p(v uint64) *uint64    { return &v }
func float32p(v float32) *float32 { return &v }

type person struct {
	FirstName string `fixed:"10"`
	LastName  string `fixed:"10"`
	Age       int    `fixed:"4"`
	Job       string `fixed:"8"`
}

type nestedStructWithTag struct {
	Cat cat `fixed:"13"`
}

type nestedStructWithoutTag struct {
	Cat cat
}

type cat struct {
	Name   string `fixed:"10"`
	Gender string `fixed:"6"`
}

type embededStruct struct {
	Number int `fixed:"3"`
	person
}

type embededStructWithTag struct {
	Number int `fixed:"3"`
	person `fixed:"15"`
}

type mixedStruct struct {
	F1 string  `fixed:"5"`
	F2 *string `fixed:"3"`
	cat
	F3  float32     `fixed:"4"`
	F4  *float64    `fixed:"4"`
	F5  interface{} `fixed:"6"`
	F6  interface{} `fixed:"8"`
	F7  cat         `fixed:"10"`
	F8  *cat
	F9  *int     `fixed:"3"`
	F10 int8     `fixed:"3"`
	F11 *int8    `fixed:"3"`
	F12 int16    `fixed:"3"`
	F13 *int16   `fixed:"3"`
	F14 int32    `fixed:"3"`
	F15 *int32   `fixed:"3"`
	F16 int64    `fixed:"3"`
	F17 *int64   `fixed:"3"`
	F18 uint     `fixed:"3"`
	F19 *uint    `fixed:"3"`
	F20 uint8    `fixed:"3"`
	F21 *uint8   `fixed:"3"`
	F22 uint16   `fixed:"3"`
	F23 *uint16  `fixed:"3"`
	F24 uint32   `fixed:"3"`
	F25 *uint32  `fixed:"3"`
	F26 uint64   `fixed:"3"`
	F27 *uint64  `fixed:"3"`
	F28 *float32 `fixed:"5"`
	F29 float64  `fixed:"5"`
	F30 *string  `fixed:"2"`
}
