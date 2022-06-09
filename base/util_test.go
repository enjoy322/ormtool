package base

import "testing"

func TestUpperCamel(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "t1", args: args{"user"}, want: "User"},
		{name: "t2", args: args{" user"}, want: "User"},
		{name: "t3", args: args{" user "}, want: "User"},
		{name: "t4", args: args{"user"}, want: "User"},
		{name: "t5", args: args{"userName"}, want: "UserName"},
		{name: "t6", args: args{"user_name"}, want: "UserName"},
		{name: "t7", args: args{""}, want: ""},
		{name: "t8", args: args{" "}, want: ""},
		{name: "t9", args: args{"  "}, want: ""},
		{name: "t9", args: args{"  user"}, want: "User"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UpperCamel(tt.args.s); got != tt.want {
				t.Errorf("UpperCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}
