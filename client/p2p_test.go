package main

import "testing"

func Test_isNeed(t *testing.T) {
	type args struct {
		f    string
		file string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				f:    "h1.txt",
				file: "h2.txt",
			},
			want: true,
		},
		{
			name: "2",
			args: args{
				f:    "hyj1.txt",
				file: "h2.txt",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNeed(tt.args.f, tt.args.file); got != tt.want {
				t.Errorf("isNeed() = %v, want %v", got, tt.want)
			}
		})
	}
}
