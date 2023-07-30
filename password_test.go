package password

import "testing"

func TestPasswordHash(t *testing.T) {
	type args struct {
		pwd string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "1",
			args:    args{pwd: "qwerpvp123"},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PasswordHash(tt.args.pwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("PasswordHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PasswordHash() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswordVerify(t *testing.T) {
	type args struct {
		pwd  string
		hash string
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
				pwd:  "qwerpvp123",
				hash: "$2a$10$y3AiH/a7mbEvxFIkbJExQe/Iz8a5rSGIfghFMmYZvFf/tKTN0nACq",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PasswordVerify(tt.args.pwd, tt.args.hash); got != tt.want {
				t.Errorf("PasswordVerify() = %v, want %v", got, tt.want)
			}
		})
	}
}
