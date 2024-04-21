package logger

import (
	"fmt"
	"testing"
)

func TestInitialize(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestInitialize", args: args{level: "debug"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Initialize(tt.args.level); (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func ExampleInitialize() {
	if err := Initialize("debug"); err != nil {
		fmt.Printf("Initialize() error = %v", err)
		return
	}
	fmt.Printf("log level %v", Log.Level())

	// Output:
	// log level debug
}
