package helper

import (
	"fmt"
	"reflect"
	"testing"
	// "unicode"
)

func TestCharEncoding(t *testing.T) {
	enc := CharEncoding("UTF1")
	fmt.Printf("%+v %[1]T\n", enc)
}

func TestFileExists(t *testing.T) {
	defer testerConfig.RecoverPanic(t)
	testerConfig.InitTest(t, "helper")

	tp := testerConfig.TempSourcePathGeter("helper")
	// Prepare test pairs
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"file_exits", args{tp("some_file.txt")}, true, false},
		{"file_not_exits", args{tp("nonexisten_file.txt")}, false, false},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FileExists(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"FileExists() error = %v, wantErr %v, args %v",
					err, tt.wantErr, tt.args)
				return
			}
			if got != tt.want {
				t.Errorf(
					"FileExists() = %v, want %v, args %v",
					got, tt.want, tt.args)
			}
		})
	}
}

func TestFilePathEncoding(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		// want    encoding.Encoding
		want    CharEncoding
		wantErr bool
	}{
		{"mini_UTF8",
			args{"/tmp/test2/RD_00-05_Radiožurnál_Friday_W06_2024_02_09.xml"},
			CharEncodingUTF8, false},
		{"mini_UTF16le",
			args{"/tmp/test2/RD_20-24_RŽ_Sport_-_út__23_07_2024_2_21282780_20240724001003.xml"},
			CharEncodingUTF16le, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FilePathEncoding(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileEncoding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFileEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileReadAllHandleEncoding(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		// want    encoding.Encoding
		want    CharEncoding
		wantErr bool
	}{
		{"mini_UTF8",
			args{"/tmp/test2/RD_00-05_Radiožurnál_Friday_W06_2024_02_09.xml"},
			CharEncodingUTF8, false},
		{"mini_UTF16le",
			args{"/tmp/test2/RD_20-24_RŽ_Sport_-_út__23_07_2024_2_21282780_20240724001003.xml"},
			CharEncodingUTF16le, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, che, err := FileReadAllHandleEncoding(tt.args.filePath)
			if err != nil {
				t.Error(err)
			}
			fmt.Println(che, string(data[0:400]))
		})
	}
}
