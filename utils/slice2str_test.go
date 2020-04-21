package utils

import "testing"

func TestSlice2Str(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "path_exist",
			args:args{s:[]string{"a","b","c"}},
			want:"a,b,c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slice2Str(tt.args.s); got != tt.want {
				t.Errorf("Slice2Str() = %v, want %v", got, tt.want)
			}
		})
	}
}