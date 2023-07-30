package mymd5

import (
	"reflect"
	"testing"
)

func TestCalcMd5_2(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "1",
			args:    args{fileName: "file-0.zip"},
			want:    []byte{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalcMd5_2(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalcMd5_2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalcMd5_2() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMd5check2(t *testing.T) {
	type args struct {
		fileName string
		md5bytes []byte
	}
	tests := []struct {
		name  string
		args  args
		want  error
		want1 bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				fileName: "file-0.zip",
				md5bytes: []byte{43, 109, 17, 140, 20, 161, 32, 235, 207, 56, 23, 25, 193, 149, 25, 13},
			},
			want:  nil,
			want1: true,
		},
		{
			name: "2",
			args: args{
				fileName: "D:\\GoProjects\\src\\p2p\\client\\file-0.zip",
				md5bytes: []byte{43, 109, 17, 14, 20, 161, 32, 235, 207, 56, 23, 25, 193, 149, 25, 13},
			},
			want:  nil,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Md5check2(tt.args.fileName, tt.args.md5bytes)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Md5check2() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Md5check2() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
