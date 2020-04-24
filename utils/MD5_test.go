package utils

import "testing"

func TestStrMd5(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "MD5Test",
			args:args{str:"crypted data"},
			want:"121d0611200c43f67b722446a4faea45",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrMd5(tt.args.str); got != tt.want {
				t.Errorf("StrMd5() = %v, want %v", got, tt.want)
			}
		})
	}
}
