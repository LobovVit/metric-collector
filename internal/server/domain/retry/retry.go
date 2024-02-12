package retry

import (
	"errors"
	"os"
	"time"

	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type Retry struct {
	attemptsLeft int
	interval     time.Duration
}

func New(attempts int) *Retry {
	ret := &Retry{attemptsLeft: attempts}
	ret.addIteration(0)
	return ret
}

func (re *Retry) runTimer() {
	time.Sleep(re.interval)
}

func (re *Retry) addIteration(delta int) {
	re.attemptsLeft = re.attemptsLeft - delta
	switch re.attemptsLeft {
	case 3:
		re.interval = 1 * time.Second
	case 2:
		re.interval = 3 * time.Second
	case 1:
		re.interval = 5 * time.Second
	default:
		re.interval = 1 * time.Second
	}
}

func (re *Retry) isRetryable(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
		return true
	}
	var osErr *os.SyscallError
	return errors.As(err, &osErr)
}

func (re *Retry) RunNoParam(f func() error) error {
	var err error
	for i := re.attemptsLeft; i > 0; i-- {
		re.addIteration(1)
		err = f()
		if !re.isRetryable(err) {
			return err
		}
		re.runTimer()
	}
	return err
}

func (re *Retry) RunMetricsParam(f func(metrics []metrics.Metrics) error, param []metrics.Metrics) error {
	var err error
	for i := re.attemptsLeft; i > 0; i-- {
		re.addIteration(1)
		err = f(param)
		if !re.isRetryable(err) {
			return err
		}
		re.runTimer()
	}
	return err
}

func (re *Retry) RunKVFloatParam(f func(key string, val float64) error, paramKey string, paramVal float64) error {
	var err error
	for i := re.attemptsLeft; i > 0; i-- {
		re.addIteration(1)
		err = f(paramKey, paramVal)
		if !re.isRetryable(err) {
			return err
		}
		re.runTimer()
	}
	return err
}
func (re *Retry) RunKVIntParam(f func(key string, val int64) error, paramKey string, paramVal int64) error {
	var err error
	for i := re.attemptsLeft; i > 0; i-- {
		re.addIteration(1)
		err = f(paramKey, paramVal)
		if !re.isRetryable(err) {
			return err
		}
		re.runTimer()
	}
	return err
}
