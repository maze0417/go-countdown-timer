package go_countdown_timer

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Timer struct {
	duration    time.Duration
	currentTick int32
	timer       *time.Timer
	ticker      *time.Ticker
	mu          sync.Mutex
	processChan chan interface{}
	stoppedChan chan interface{}
}

func NewTimer(duration time.Duration) *Timer {
	return &Timer{
		duration:    duration,
		stoppedChan: make(chan interface{}, 1),
		processChan: make(chan interface{}),
	}
}

func (c *Timer) start(processFunc func(interface{}) bool, timeoutFunc func()) {

	c.ticker = time.NewTicker(time.Second)
	c.timer = time.NewTimer(c.duration)
	v := int32(c.duration.Seconds())

	atomic.StoreInt32(&c.currentTick, v)
	c.counting(c.ticker, timeoutFunc, processFunc)
}

func (c *Timer) StartWithActions(processFunc func(interface{}) bool, timeoutFunc func()) {
	c.start(processFunc, timeoutFunc)
}

func (c *Timer) resetTick() {
	if c.ticker != nil {
		c.ticker.Stop()
		c.ticker = nil
	}

	//clear chan
	if len(c.processChan) > 0 {
		<-c.processChan
	}

	if len(c.stoppedChan) > 0 {
		<-c.stoppedChan
	}

	t := atomic.LoadInt32(&c.currentTick)
	if t > 0 {
		atomic.StoreInt32(&c.currentTick, 0)
	}
	c.stoppedChan <- 0
}

func (c *Timer) Stop(from ...string) {
	<-c.stoppedChan
}

func (c *Timer) GetProcessChan() chan<- interface{} {
	return c.processChan
}

func (c *Timer) counting(t *time.Ticker, timeoutProcessFunc func(), processFunc func(interface{}) bool) {
	defer c.resetTick()

	for {
		v := atomic.LoadInt32(&c.currentTick)

		select {
		case <-t.C:
			fmt.Printf("Counting down %v \n", v)
			v--
			atomic.StoreInt32(&c.currentTick, v)
			if v == 0 {
				if timeoutProcessFunc != nil {
					timeoutProcessFunc()
				}
				return
			}

		case arg := <-c.processChan:
			if processFunc == nil {
				continue
			}

			if processFunc(arg) {
				return
			}

		case <-c.timer.C:
			return
		}

	}
}
