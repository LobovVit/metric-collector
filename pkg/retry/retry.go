package retry

import (
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

func (re *Retry) runTimer() {
	time.Sleep(re.interval)
}

func (re *Retry) addIteration(delta int) {
	re.interval = re.interval + time.Duration(re.curAttempts+1)*time.Second
	re.curAttempts = re.curAttempts + delta
}

func (re *Retry) Run() bool {
	re.addIteration(1)
	if re.maxAttempts < re.curAttempts {
		return true
	}
	logger.Log.Info("Retry", zap.String("current attempt", re.String()))
	re.runTimer()
	return false
}

func (re *Retry) String() string {
	return "attempt=" + strconv.Itoa(re.curAttempts) + " delay=" + re.interval.String()
}
