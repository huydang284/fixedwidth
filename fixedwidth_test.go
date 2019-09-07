package fixedwidth

import (
	"fmt"
	"reflect"
	"testing"
)

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

func TestMarshal(t *testing.T) {
	type args struct {
		v interface{}
	}

	p := person{
		FirstName: "Alexander",
		LastName:  "Goodword",
		Age:       40,
		Job:       "Software Engineer",
	}
	p2 := []person{
		p,
		{
			FirstName: "Frank",
			LastName:  "Lampard",
			Age:       41,
			Job:       "Coach",
		},
	}
	p3 := person{
		FirstName: "นายทดสอบ",
		LastName:  "ทดสอบ",
		Age:       100,
		Job:       "Retired",
	}
	tests := []struct {
		name    string
		args    args
		want    []rune
		wantErr bool
	}{
		{
			name:    "single line",
			args:    args{v: p},
			want:    []rune("Alexander Goodword  40  Software"),
			wantErr: false,
		},
		{
			name:    "double lines",
			args:    args{v: p2},
			want:    []rune("Alexander Goodword  40  Software\nFrank     Lampard   41  Coach   "),
			wantErr: false,
		},
		{
			name:    "single line - unicode",
			args:    args{v: p3},
			want:    []rune("นายทดสอบ  ทดสอบ     100 Retired "),
			wantErr: false,
		},
		{
			name:    "nested struct with tag",
			args:    args{v: nestedStructWithTag{Cat: cat{Name: "June", Gender: "male"}}},
			want:    []rune("June      mal"),
			wantErr: false,
		},
		{
			name:    "nested struct without tag",
			args:    args{v: nestedStructWithoutTag{Cat: cat{Name: "June", Gender: "male"}}},
			want:    []rune("June      male  "),
			wantErr: false,
		},
		{
			name: "embeded struct without tag",
			args: args{v: embededStruct{
				Number: 15,
				person: person{
					FirstName: "Drogba",
					LastName:  "Didier",
					Age:       41,
					Job:       "Retired",
				},
			}},
			want:    []rune("15 Drogba    Didier    41  Retired "),
			wantErr: false,
		},
		{
			name: "embeded struct with tag",
			args: args{v: embededStructWithTag{
				Number: 15,
				person: person{
					FirstName: "Drogba",
					LastName:  "Didier",
					Age:       41,
					Job:       "Retired",
				},
			}},
			want:    []rune("15 Drogba    Didie"),
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
			fmt.Println(string(got))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Marshal() got = %v, want %v", got, tt.want)
			}
		})
	}
}
