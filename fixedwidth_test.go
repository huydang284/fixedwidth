package fixedwidth

func float64p(v float64) *float64 { return &v }
func stringp(v string) *string    { return &v }

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
	F3 float32     `fixed:"4"`
	F4 *float64    `fixed:"4"`
	F5 interface{} `fixed:"6"`
	F6 interface{} `fixed:"8"`
	F7 cat         `fixed:"10"`
	F8 *cat
}
