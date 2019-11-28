package fixedwidth

import "testing"

var s = mixedStructForUnmarshal{
	F1: "the f",
	F2: stringp("sec"),
	cat: cat{
		Name:   "P",
		Gender: "female",
	},
	F3: 10.5,
	F4: float64p(7.22),
	F5: "what i",
	F6: "7",
	F7: cat{
		Name:   "Ali",
		Gender: "",
	},
	F8: &cat{
		Name:   "wow",
		Gender: "male",
	},
	F9:  intp(1),
	F10: 2,
	F11: int8p(3),
	F12: 4,
	F13: int16p(5),
	F14: 6,
	F15: int32p(7),
	F16: 8,
	F17: int64p(9),
	F18: 1,
	F19: uintp(2),
	F20: 3,
	F21: uint8p(4),
	F22: 5,
	F23: uint16p(6),
	F24: 7,
	F25: uint32p(8),
	F26: 9,
	F27: uint64p(10),
	F28: float32p(1.12),
	F29: 2.23,
	F30: stringp(""),
}

func BenchmarkMarshal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Marshal(s)
	}
}

type mixedData struct {
	F1  string   `fixed:"1,10"`
	F2  *string  `fixed:"11,20"`
	F3  int64    `fixed:"21,30"`
	F4  *int64   `fixed:"31,40"`
	F5  int32    `fixed:"41,50"`
	F6  *int32   `fixed:"51,60"`
	F7  int16    `fixed:"61,70"`
	F8  *int16   `fixed:"71,80"`
	F9  int8     `fixed:"81,90"`
	F10 *int8    `fixed:"91,100"`
	F11 float64  `fixed:"101,110"`
	F12 *float64 `fixed:"111,120"`
	F13 float32  `fixed:"121,130"`
}

var mixedDataInstance = mixedData{"foo", stringp("foo"), 42, int64p(42), 42, int32p(42), 42, int16p(42), 42, int8p(42), 4.2, float64p(4.2), 4.2} //,float32p(4.2)}

func BenchmarkMarshal_MixedData_1000(b *testing.B) {
	v := make([]mixedData, 1000)
	for i := range v {
		v[i] = mixedDataInstance
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Marshal(v)
	}
}
