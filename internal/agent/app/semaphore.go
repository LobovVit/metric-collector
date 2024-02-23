package app

type semaphore struct {
	semaCh chan struct{}
}

func newSemaphore(maxReq int64) *semaphore {
	return &semaphore{
		semaCh: make(chan struct{}, maxReq),
	}
}

func (s *semaphore) acquire() {
	s.semaCh <- struct{}{}
}

func (s *semaphore) release() {
	<-s.semaCh
}
