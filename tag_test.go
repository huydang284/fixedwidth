package fixedwidth

import (
	"reflect"
	"testing"
)

func Test_tag_getLimitFixedTag(t *testing.T) {
	type validFixedTag struct {
		Name string `fixed:"10"`
	}
	var v validFixedTag
	type invalidFixedTag struct {
		Name string `fixed:"abc"`
	}
	var i invalidFixedTag

	type args struct {
		field reflect.StructField
	}
	tests := []struct {
		name string
		args args
		want int
		ok   bool
	}{
		{
			name: "success",
			args: args{
				field: reflect.TypeOf(v).Field(0),
			},
			want: 10,
			ok:   true,
		},
		{
			name: "fail",
			args: args{
				field: reflect.TypeOf(i).Field(0),
			},
			want: 0,
			ok:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := tag{}
			got, got1 := ta.getLimitFixedTag(tt.args.field)
			if got != tt.want {
				t.Errorf("getLimitFixedTag() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.ok {
				t.Errorf("getLimitFixedTag() got1 = %v, want %v", got1, tt.ok)
			}
		})
	}
}
