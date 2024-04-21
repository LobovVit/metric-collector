// Package retry - included retry functions
package retry

import (
	"context"
	"time"
)

// Do - retry function with one any parameter and context, function returning only error
func Do[T any](ctx context.Context, repeat int, retryFunc func(context.Context, T) error, p T, isRepeatableFunc func(err error) bool) error {
	var err error
	for i := 0; i < repeat; i++ {
		// Return immediately if ctx is canceled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = retryFunc(ctx, p)
		if err == nil || !isRepeatableFunc(err) {
			break
		}
		if i < repeat-1 {
			time.Sleep(time.Second * 2 * time.Duration(i+1))
		}
	}
	return err
}

// DoTwoParams - retry function with two any parameter and context, function returning only error
func DoTwoParams[T1, T2 any](ctx context.Context, repeat int, retryFunc func(context.Context, T1, T2) error, p1 T1, p2 T2, isRepeatableFunc func(err error) bool) error {
	var err error
	for i := 0; i < repeat; i++ {
		// Return immediately if ctx is canceled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = retryFunc(ctx, p1, p2)
		if err == nil || !isRepeatableFunc(err) {
			break
		}
		if i < repeat-1 {
			time.Sleep(time.Second * 2 * time.Duration(i+1))
		}
	}
	return err
}

// DoNoParams - retry function without any parameter only context, function returning only error
func DoNoParams(ctx context.Context, repeat int, retryFunc func(context.Context) error, isRepeatableFunc func(err error) bool) error {
	var err error
	for i := 0; i < repeat; i++ {
		// Return immediately if ctx is canceled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = retryFunc(ctx)
		if err == nil || !isRepeatableFunc(err) {
			break
		}
		if i < repeat-1 {
			time.Sleep(time.Second * 2 * time.Duration(i+1))
		}
	}
	return err
}
