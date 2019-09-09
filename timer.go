package go_countdown_timer

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Timer struct {
	duration     time.Duration
	currentTick  int32
	timer        *time.Timer
	ticker       *time.Ticker
	processChan  chan func() bool
	stoppingChan chan interface{}
	stoppedChan  chan interface{}
	mu           sync.Mutex
}

func NewTimer(duration time.Duration) *Timer {
	return &Timer{
		duration:     duration,
		stoppedChan:  make(chan interface{}, 1),
		stoppingChan: make(chan interface{}, 1),
		processChan:  make(chan func() bool),
	}
}

func (c *Timer) StartNew(timeoutFunc func(), from ...string) {
	c.mu.Lock()
	c.ensurePreviousStopIfAny(from...)
	c.ticker = time.NewTicker(time.Second)
	c.timer = time.NewTimer(c.duration)
	v := int32(c.duration.Seconds())
	c.mu.Unlock()

	atomic.StoreInt32(&c.currentTick, v)
	c.countdown(c.ticker, timeoutFunc, from...)
}

func (c *Timer) resetTick() {
	if c.ticker != nil {
		c.ticker.Stop()
		c.ticker = nil
	}
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
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

func (c *Timer) ensurePreviousStopIfAny(from ...string) {
	if c.ticker == nil {
		fmt.Printf("timer %v is new \n", from)
		return
	}
	fmt.Printf("timer %v is stopping another timer \n", from)
	c.stoppingChan <- 0

	<-c.stoppedChan
	fmt.Printf("timer %v stop other done \n", from)
}
func (c *Timer) ReceiveProcessEvent() chan<- func() bool {
	return c.processChan
}

func (c *Timer) countdown(t *time.Ticker, timeoutProcessFunc func(), from ...string) {
	defer c.resetTick()
	for {
		v := atomic.LoadInt32(&c.currentTick)

		select {
		case <-t.C:
			fmt.Printf("timer %v Counting down %v \n", from, v)
			v--
			atomic.StoreInt32(&c.currentTick, v)
			if v == 0 {
				if timeoutProcessFunc != nil {
					timeoutProcessFunc()
				}
				return
			}

		case arg := <-c.processChan:

			if arg == nil {
				continue
			}

			if arg() {
				return
			}
		case <-c.stoppingChan:
			return
		case <-c.timer.C:
			return
		}

	}
}
