package fixedwidth

import (
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"log"
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		data string
		want interface{}
	}{
		{
			name: "single line ",
			data: "Huy       Đặng      25  Engineer",
			want: []person{
				{
					FirstName: "Huy",
					LastName:  "Đặng",
					Age:       25,
					Job:       "Engineer",
				},
			},
		},
		{
			name: "multiple lines",
			data: "Huy       Dang      25  Engineer\nDidier    Drogba    41  Retired \nLâm       Đặng      26  Thủ môn ",
			want: []person{
				{
					FirstName: "Huy",
					LastName:  "Dang",
					Age:       25,
					Job:       "Engineer",
				},
				{
					FirstName: "Didier",
					LastName:  "Drogba",
					Age:       41,
					Job:       "Retired",
				},
				{
					FirstName: "Lâm",
					LastName:  "Đặng",
					Age:       26,
					Job:       "Thủ môn",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p []person
			err := Unmarshal([]byte(tt.data), &p)
			if err != nil {
				t.Error(err)
				return
			}

			if !reflect.DeepEqual(p, tt.want) {
				t.Error(errors.New("incorrect result"))
			}
		})
	}

	t.Run("nested struct with tag", func(t *testing.T) {
		want := nestedStructWithTag{Cat: cat{Name: "June", Gender: "mal"}}
		var s nestedStructWithTag
		err := Unmarshal([]byte("June      mal"), &s)
		if err != nil {
			t.Error(err)
			return
		}

		if !reflect.DeepEqual(s, want) {
			t.Error(errors.New("incorrect result"))
		}
	})

	t.Run("nested struct without tag", func(t *testing.T) {
		want := nestedStructWithoutTag{Cat: cat{Name: "June", Gender: "male"}}
		var s nestedStructWithoutTag
		err := Unmarshal([]byte("June      male  "), &s)
		if err != nil {
			t.Error(err)
			return
		}

		if !reflect.DeepEqual(s, want) {
			t.Error(errors.New("incorrect result"))
		}
	})

	t.Run("embeded struct with tag", func(t *testing.T) {
		want := embededStructWithTag{
			Number: 15,
			person: person{
				FirstName: "Drogba",
				LastName:  "Didie",
				Age:       0,
				Job:       "",
			},
		}
		var s embededStructWithTag
		err := Unmarshal([]byte("15 Drogba    Didie"), &s)
		if err != nil {
			t.Error(err)
			return
		}

		if !reflect.DeepEqual(s, want) {
			t.Error(errors.New("incorrect result"))
		}
	})

	t.Run("embeded struct without tag", func(t *testing.T) {
		want := embededStruct{
			Number: 15,
			person: person{
				FirstName: "Drogba",
				LastName:  "Didier",
				Age:       41,
				Job:       "Retired",
			},
		}
		var s embededStruct
		err := Unmarshal([]byte("15 Drogba    Didier    41  Retired "), &s)
		if err != nil {
			t.Error(err)
			return
		}

		if !reflect.DeepEqual(s, want) {
			t.Error(errors.New("incorrect result"))
		}
	})

	t.Run("mixed type", func(t *testing.T) {
		want := mixedStructForUnmarshal{
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
		var s mixedStruct
		err := Unmarshal([]byte("the fsecP         female10.57.22what i7       Ali       wow       male  1  2  3  4  5  6  7  8  9  1  2  3  4  5  6  7  8  9  10 1.12 2.23   "), &s)
		if err != nil {
			t.Error(err)
			return
		}

		if !cmp.Equal(s.F1, want.F1) {
			t.Error("incorrect F1")
			return
		}
		if !cmp.Equal(s.F2, want.F2) {
			t.Error("incorrect F2")
			return
		}
		if !cmp.Equal(s.cat, want.cat) {
			t.Error("incorrect cat")
			return
		}
		if !cmp.Equal(s.F3, want.F3) {
			t.Error("incorrect F3")
			return
		}
		if !cmp.Equal(s.F4, want.F4) {
			t.Error("incorrect F4")
			return
		}
		if !cmp.Equal(s.F5, want.F5) {
			t.Error("incorrect F5")
			return
		}
		if !cmp.Equal(s.F6, want.F6) {
			t.Error("incorrect F6")
			return
		}
		if !cmp.Equal(s.F7, want.F7) {
			t.Error("incorrect F7")
			return
		}
		if !cmp.Equal(s.F8, want.F8) {
			t.Error("incorrect F8")
			return
		}
		if !cmp.Equal(s.F9, want.F9) {
			t.Error("incorrect F9")
			return
		}
		if !cmp.Equal(s.F10, want.F10) {
			t.Error("incorrect F10")
			return
		}
		if !cmp.Equal(s.F11, want.F11) {
			t.Error("incorrect F11")
			return
		}
		if !cmp.Equal(s.F12, want.F12) {
			t.Error("incorrect F12")
			return
		}
		if !cmp.Equal(s.F13, want.F13) {
			t.Error("incorrect F13")
			return
		}
		if !cmp.Equal(s.F14, want.F14) {
			t.Error("incorrect F14")
			return
		}
		if !cmp.Equal(s.F15, want.F15) {
			t.Error("incorrect F15")
			return
		}
		if !cmp.Equal(s.F16, want.F16) {
			t.Error("incorrect F16")
			return
		}
		if !cmp.Equal(s.F17, want.F17) {
			t.Error("incorrect F17")
			return
		}
		if !cmp.Equal(s.F18, want.F18) {
			t.Error("incorrect F18")
			return
		}
		if !cmp.Equal(s.F19, want.F19) {
			t.Error("incorrect F19")
			return
		}
		if !cmp.Equal(s.F20, want.F20) {
			t.Error("incorrect F20")
			return
		}
		if !cmp.Equal(s.F21, want.F21) {
			t.Error("incorrect F21")
			return
		}
		if !cmp.Equal(s.F22, want.F22) {
			t.Error("incorrect F22")
			return
		}
		if !cmp.Equal(s.F23, want.F23) {
			t.Error("incorrect F23")
			return
		}
		if !cmp.Equal(s.F24, want.F24) {
			t.Error("incorrect F24")
			return
		}
		if !cmp.Equal(s.F25, want.F25) {
			t.Error("incorrect F25")
			return
		}
		if !cmp.Equal(s.F26, want.F26) {
			t.Error("incorrect F26")
			return
		}
		if !cmp.Equal(s.F27, want.F27) {
			t.Error("incorrect F27")
			return
		}
		if !cmp.Equal(s.F28, want.F28) {
			t.Error("incorrect F28")
			return
		}
		if !cmp.Equal(s.F29, want.F29) {
			t.Error("incorrect F29")
			return
		}
		if !cmp.Equal(s.F30, want.F30) {
			t.Error("incorrect F30")
			return
		}
	})
}

func ExampleUnmarshaler_Unmarshal() {
	var p person
	m := NewUnmarshaler()
	err := m.Unmarshal([]byte("Alexander Goodword  40  Software"), &p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", p)
	// Output:
	// {FirstName:Alexander LastName:Goodword Age:40 Job:Software}
}
