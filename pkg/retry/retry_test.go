package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

var retryable = errors.New("retryable")

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, retryable)
}

func oneParamR(ctx context.Context, p1 int) error {
	p1++
	return retryable
}

func oneParamN(ctx context.Context, p1 int) error {
	p1++
	return nil
}

func TestDo(t *testing.T) {
	type args struct {
		ctx              context.Context
		repeat           int
		retryFunc        func(context.Context, int) error
		p                int
		isRepeatableFunc func(err error) bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test1", args: args{ctx: context.Background(), repeat: 3, retryFunc: oneParamR, p: 1, isRepeatableFunc: IsRetryable}, wantErr: true},
		{name: "test2", args: args{ctx: context.Background(), repeat: 3, retryFunc: oneParamN, p: 1, isRepeatableFunc: IsRetryable}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Do(tt.args.ctx, tt.args.repeat, tt.args.retryFunc, tt.args.p, tt.args.isRepeatableFunc); (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func noParamR(ctx context.Context) error {
	var i int = 0
	i++
	return retryable
}

func noParamN(ctx context.Context) error {
	var i int = 0
	i++
	return nil
}

func TestDoNoParams(t *testing.T) {
	type args struct {
		ctx              context.Context
		repeat           int
		retryFunc        func(context.Context) error
		isRepeatableFunc func(err error) bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test1", args: args{ctx: context.Background(), repeat: 3, retryFunc: noParamR, isRepeatableFunc: IsRetryable}, wantErr: true},
		{name: "test2", args: args{ctx: context.Background(), repeat: 3, retryFunc: noParamN, isRepeatableFunc: IsRetryable}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DoNoParams(tt.args.ctx, tt.args.repeat, tt.args.retryFunc, tt.args.isRepeatableFunc); (err != nil) != tt.wantErr {
				t.Errorf("DoNoParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func twoParamR(ctx context.Context, p1, p2 int) error {
	p1++
	p2++
	return retryable
}

func twoParamN(ctx context.Context, p1, p2 int) error {
	p1++
	p2++
	return nil
}

func TestDoTwoParams(t *testing.T) {
	type args struct {
		ctx              context.Context
		repeat           int
		retryFunc        func(context.Context, int, int) error
		p1               int
		p2               int
		isRepeatableFunc func(err error) bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test1", args: args{ctx: context.Background(), repeat: 3, retryFunc: twoParamR, p1: 1, p2: 1, isRepeatableFunc: IsRetryable}, wantErr: true},
		{name: "test1", args: args{ctx: context.Background(), repeat: 3, retryFunc: twoParamN, p1: 1, p2: 1, isRepeatableFunc: IsRetryable}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DoTwoParams(tt.args.ctx, tt.args.repeat, tt.args.retryFunc, tt.args.p1, tt.args.p2, tt.args.isRepeatableFunc); (err != nil) != tt.wantErr {
				t.Errorf("DoTwoParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func ExampleDo() {
	err := Do(context.Background(), 3, oneParamR, 1, IsRetryable)
	if err != nil {
		fmt.Printf(err.Error())
	}

	// Output:
	// retryable
}
