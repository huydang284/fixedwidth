package fixedwidth

import (
	"errors"
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
			name: "single line",
			data: "Huy       Dang      25  Engineer",
			want: person{
				FirstName: "Huy",
				LastName:  "Dang",
				Age:       25,
				Job:       "Engineer",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p person
			err := Unmarshal([]rune(tt.data), &p)
			if err != nil {
				t.Error(err)
				return
			}

			if !reflect.DeepEqual(p, tt.want) {
				t.Error(errors.New("incorrect result"))
			}
		})
	}
}
