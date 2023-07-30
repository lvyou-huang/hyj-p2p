package processfile

import "testing"

func TestMakeSilce(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{fileName: "C:\\Users\\86188\\Videos\\2022211940-hyj.mp4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MakeSilce(tt.args.fileName)
		})
	}
}

func TestMakeFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "1",
			args:    args{filePath: "C:\\Users\\86188\\Desktop\\c\\file.zip"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MakeFile(tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("MakeFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMakeSilce1(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{filePath: "C:\\Users\\86188\\Desktop\\c\\c.zip"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MakeSilce(tt.args.filePath)
		})
	}
}
