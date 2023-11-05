package testhelpers

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrWaitGroupTimeout = errors.New("timeout waiting for WaitGroup")
)

func WaitOrTimeout(wg *sync.WaitGroup, timeout time.Duration) error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return nil
	case <-time.After(timeout):
		return ErrWaitGroupTimeout
	}
}
