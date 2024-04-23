package compress

import (
	"reflect"
	"testing"
)

func TestCompress(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		arg     args
		wantErr bool
	}{
		{name: "test compress #1", arg: args{data: []byte("hello world")}, wantErr: false},
		{name: "test compress #2", arg: args{data: []byte("hello \n world")}, wantErr: false},
		{name: "test compress #3", arg: args{data: []byte("{hello world}")}, wantErr: false},
		{name: "test compress #4", arg: args{data: []byte("Привет мир")}, wantErr: false},
		{name: "test compress #5", arg: args{data: []byte("!\"№%:,.;())}")}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compress(tt.arg.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.DeepEqual(got, tt.arg.data) {
				t.Errorf("not compressed got = %v, arg %v", got, tt.arg.data)
			}
			ungot, err := UnCompress(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnCompress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.arg.data, ungot) {
				t.Errorf("Compress() got = %v, UnCompress() got = %v", got, ungot)
			}
		})
	}
}
