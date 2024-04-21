package signature

import (
	"fmt"
	"testing"
)

func TestSignature(t *testing.T) {
	type args struct {
		data []byte
		key  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "valid-1",
			args:    args{data: []byte("test"), key: "sec-key"},
			wantErr: false,
		}, {
			name:    "valid-2",
			args:    args{data: []byte("te\nst"), key: "sec\nkey"},
			wantErr: false,
		}, {
			name:    "valid-3",
			args:    args{data: []byte("test"), key: ""},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := CreateSignature(tt.args.data, tt.args.key)
			if err != nil {
				t.Errorf("CreateSignature() error = %v", err)
			}
			if hash != nil {
				err = CheckSignature(tt.args.data, fmt.Sprintf("%x", hash), tt.args.key)
				if err != nil && !tt.wantErr {
					t.Errorf("CreateSignature() error = %v", err)
				}
			}

		})
	}
}

func ExampleCreateSignature() {
	data := []byte("test")
	key := "sec-key"
	ret, err := CreateSignature(data, key)
	if err != nil {
		fmt.Printf("CreateSignature() error = %v", err)
	}
	fmt.Println(ret)

	// Output:
	// [237 39 112 44 62 34 161 216 18 133 22 130 227 234 236 143 215 102 102 235 236 196 224 157 237 210 245 158 55 97 172 163]
}

func ExampleCheckSignature() {
	data := []byte("test")
	key := "sec-key"
	hash := []byte{237, 39, 112, 44, 62, 34, 161, 216, 18, 133, 22, 130, 227, 234,
		236, 143, 215, 102, 102, 235, 236, 196, 224, 157, 237, 210, 245, 158, 55, 97, 172, 163}
	err := CheckSignature(data, fmt.Sprintf("%x", hash), key)
	if err != nil {
		fmt.Printf("CreateSignature() error = %v", err)
		return
	}
	fmt.Printf("CreateSignature %v", "OK")

	// Output:
	// CreateSignature OK
}
