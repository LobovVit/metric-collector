package retry

import (
	"context"
	"strconv"
	"time"

	"github.com/LobovVit/metric-collector/pkg/logger"
	"go.uber.org/zap"
)

type Retry struct {
	curAttempts int
	maxAttempts int
	interval    time.Duration
}

func New(attempts int) *Retry {
	ret := &Retry{maxAttempts: attempts}
	return ret
}

func (re *Retry) sleep() {
	time.Sleep(re.interval)
}

func (re *Retry) addIteration(delta int) {
	re.interval = re.interval + time.Duration(re.curAttempts+1)*time.Second
	re.curAttempts = re.curAttempts + delta
}

func (re *Retry) Run() bool {
	re.addIteration(1)
	if re.maxAttempts < re.curAttempts {
		return false
	}
	logger.Log.Info("Retry", zap.String("current attempt", re.String()))
	re.sleep()
	return true
}

func (re *Retry) String() string {
	return "attempt=" + strconv.Itoa(re.curAttempts) + " delay=" + re.interval.String()
}

func DoWithoutReturn[T any](ctx context.Context, repeat int, retryFunc func(context.Context, T) error, p T, isRepeatableFunc func(err error) bool) error {
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
			time.Sleep(time.Second * 3)
		}
	}
	return err
}

func DoWithReturn[T, R any](ctx context.Context, repeat int, retryFunc func(context.Context, T) (R, error), p T, isRepeatableFunc func(err error) bool) (R, error) {
	var err error
	var ret R
	for i := 0; i < repeat; i++ {
		// Return immediately if ctx is canceled
		select {
		case <-ctx.Done():
			return ret, ctx.Err()
		default:
		}

		ret, err = retryFunc(ctx, p)
		if err == nil || !isRepeatableFunc(err) {
			break
		}
		if i < repeat-1 {
			time.Sleep(time.Second * 3)
		}
	}
	return ret, err
}

func DoWithoutReturnNoParams(ctx context.Context, repeat int, retryFunc func(context.Context) error, isRepeatableFunc func(err error) bool) error {
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
			time.Sleep(time.Second * 3)
		}
	}
	return err
}
