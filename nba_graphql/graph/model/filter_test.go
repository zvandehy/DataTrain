package model

import "testing"

func TestPrint(t *testing.T) {
	type args struct {
		Str            string
		StrPtr         *string
		Int            int
		IntPtr         *int
		StrSlice       []string
		StrPtrSlice    []*string
		StrSlicePtr    *[]string
		StrPtrSlicePtr *[]*string
		Obj            struct {
			//TODO: need to test unexported fields
			Foo string
		}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "string",
			args: args{
				Str: "hello world",
			},
			want: "{Str:hello world}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Print(tt.args); got != tt.want {
				t.Errorf("Print() = %v, want %v", got, tt.want)
			}
		})
	}
}
