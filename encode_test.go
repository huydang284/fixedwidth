package fixedwidth

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	type args struct {
		v interface{}
	}

	singleLine := person{
		FirstName: "Alexander",
		LastName:  "Goodword",
		Age:       40,
		Job:       "Software Engineer",
	}
	wantedSingleLine := "Alexander Goodword  40  Software"

	multiLines := []person{
		{
			FirstName: "Alexander",
			LastName:  "Goodword",
			Age:       40,
			Job:       "Software Engineer",
		},
		{
			FirstName: "Frank",
			LastName:  "Lampard",
			Age:       41,
			Job:       "Coach",
		},
		{
			FirstName: "Mason",
			LastName:  "Mount",
			Age:       20,
			Job:       "Midfielder",
		},
	}
	wantedMultiLines := "Alexander Goodword  40  Software\nFrank     Lampard   41  Coach   \nMason     Mount     20  Midfield"

	singleLineUnicode := person{
		FirstName: "Huy",
		LastName:  "Đặng",
		Age:       100,
		Job:       "Kỹ sư",
	}
	wantedSingleLineUnicode := "Huy       Đặng      100 Kỹ sư   "

	multiLinesUnicode := []person{
		{
			FirstName: "Huy",
			LastName:  "Đặng",
			Age:       25,
			Job:       "Kỹ sư",
		},
		{
			FirstName: "日本人",
			LastName:  "の氏名",
			Age:       30,
			Job:       "エンジニア",
		},
		{
			FirstName: "후이",
			LastName:  "당",
			Age:       45,
			Job:       "기사",
		},
	}
	wantedMultiLinesUnicode := "Huy       Đặng      25  Kỹ sư   \n日本人       の氏名       30  エンジニア   \n후이        당         45  기사      "

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "single line",
			args:    args{v: singleLine},
			want:    []byte(wantedSingleLine),
			wantErr: false,
		},
		{
			name:    "multi lines",
			args:    args{v: multiLines},
			want:    []byte(wantedMultiLines),
			wantErr: false,
		},
		{
			name:    "single line - unicode",
			args:    args{v: singleLineUnicode},
			want:    []byte(wantedSingleLineUnicode),
			wantErr: false,
		},
		{
			name:    "multi lines - unicode",
			args:    args{v: multiLinesUnicode},
			want:    []byte(wantedMultiLinesUnicode),
			wantErr: false,
		},
		{
			name:    "nested struct with tag",
			args:    args{v: nestedStructWithTag{Cat: cat{Name: "June", Gender: "male"}}},
			want:    []byte("June      mal"),
			wantErr: false,
		},
		{
			name:    "nested struct without tag",
			args:    args{v: nestedStructWithoutTag{Cat: cat{Name: "June", Gender: "male"}}},
			want:    []byte("June      male  "),
			wantErr: false,
		},
		{
			name: "embedded struct without tag",
			args: args{v: embeddedStruct{
				Number: 15,
				person: person{
					FirstName: "Drogba",
					LastName:  "Didier",
					Age:       41,
					Job:       "Retired",
				},
			}},
			want:    []byte("15 Drogba    Didier    41  Retired "),
			wantErr: false,
		},
		{
			name: "embedded struct with tag",
			args: args{v: embeddedStructWithTag{
				Number: 15,
				person: person{
					FirstName: "Drogba",
					LastName:  "Didier",
					Age:       41,
					Job:       "Retired",
				},
			}},
			want:    []byte("15 Drogba    Didie"),
			wantErr: false,
		},
		{
			name: "mixed types",
			args: args{v: mixedStruct{
				F1: "the first field",
				F2: stringp("second"),
				cat: cat{
					Name:   "P",
					Gender: "female",
				},
				F3: 10.544,
				F4: float64p(7.222),
				F5: "what is nil",
				F6: 7,
				F7: cat{
					Name:   "Ali",
					Gender: "male",
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
			}},
			want:    []byte("the fsecP         female10.57.22what i7       Ali       wow       male  1  2  3  4  5  6  7  8  9  1  2  3  4  5  6  7  8  9  10 1.12 2.23   "),
			wantErr: false,
		},
		{
			name:    "empty struct - single line",
			args:    args{v: person{}},
			want:    []byte("                    0           "),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Marshal() got = %v, want %v \n %s, %s", got, tt.want, string(got), string(tt.want))
			}
		})
	}
}

func ExampleMarshaler_Marshal() {
	p := person{
		FirstName: "Alexander",
		LastName:  "Goodword",
		Age:       40,
		Job:       "Software Engineer",
	}
	m := NewMarshaler()
	b, err := m.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(b))
	// Output:
	// Alexander Goodword  40  Software
}
