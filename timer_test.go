package countdown

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimerCanProcessAndStop(t *testing.T) {

	timer := NewTimer(3 * time.Second)

	countProcess := 0

	go func() {
		time.Sleep(time.Second * 1)
		timer.ReceiveProcessEvent() <- func() bool {
			countProcess++
			return false
		}
	}()

	expired := false

	timer.StartNew(func() {

		expired = true
	})
	assert.Equal(t, 1, countProcess)
	assert.True(t, expired)

}
func TestTimerCanProcessMultipleChannel(t *testing.T) {

	timer := NewTimer(3 * time.Second)

	marker := map[string]interface{}{
		"foo": nil,
		"bar": nil,
	}
	go func() {
		time.Sleep(time.Second * 2)
		timer.ReceiveProcessEvent() <- func() bool {

			id := "foo"
			if _, ok := marker[id]; !ok {
				return false
			}

			marker[id] = "foo"
			cur := 0
			for _, v := range marker {
				if v != nil {
					cur++
				}
			}

			if cur == len(marker) {
				t.Logf("collected done :%v \n", marker)

				return true
			}
			return false
		}
		timer.ReceiveProcessEvent() <- func() bool {

			id := "bar"
			if _, ok := marker[id]; !ok {
				return false
			}

			marker[id] = "bar"
			cur := 0
			for _, v := range marker {
				if v != nil {
					cur++
				}
			}

			if cur == len(marker) {
				t.Logf("collected done :%v \n", marker)

				return true
			}
			return false
		}
	}()
	timer.StartNew(func() {
		t.Logf("process something after timeout")
	})
	t.Logf("marker: %v \n", marker)

	assert.Equal(t, "foo", marker["foo"])
	assert.Equal(t, "bar", marker["bar"])

}
func TestTimerCanProcessMultipleChannelAndTimeout(t *testing.T) {

	timer := NewTimer(3 * time.Second)

	marker := map[string]interface{}{
		"foo": nil,
		"bar": nil,
	}
	go func() {
		time.Sleep(time.Second * 2)
		timer.ReceiveProcessEvent() <- func() bool {

			id := "foo"
			if _, ok := marker[id]; !ok {
				return false
			}

			marker[id] = "foo"
			cur := 0
			for _, v := range marker {
				if v != nil {
					cur++
				}
			}

			if cur == len(marker) {
				t.Logf("collected done :%v", marker)
				return true
			}
			return false
		}

	}()

	timer.StartNew(func() {
		t.Logf("process something after timeout \n")
	})
	t.Logf("marker: %v", marker)
	assert.Equal(t, "foo", marker["foo"])
	assert.Equal(t, nil, marker["bar"])
}
func TestTimerStartAtTheSameTime(t *testing.T) {

	timer := NewTimer(3 * time.Second)
	timerEndChan := make(chan bool, 1)
	timerEnd := false

	go func() {
		timer.StartNew(func() {
			t.Logf("timer 1 time up")
			timerEnd = true
			timerEndChan <- true

		}, "timer1")
	}()
	go func() {
		timer.StartNew(func() {
			t.Logf("timer 2 time up")
			timerEnd = true
			timerEndChan <- true
		}, "timer2")
	}()

	assert.True(t, <-timerEndChan)
	assert.True(t, timerEnd)
}
func TestTimerStartWhenAnotherNonStop(t *testing.T) {

	timer := NewTimer(3 * time.Second)
	timerEndChan := make(chan bool, 1)
	timerEnd := false

	go func() {
		time.Sleep(time.Second)
		timer.StartNew(func() {
			t.Logf("timer 1 time up")
			timerEnd = true
			timerEndChan <- true

		}, "timer1")
	}()

	go func() {
		timer.StartNew(func() {
			t.Logf("timer 2 time up")
			timerEnd = true
			timerEndChan <- true
		}, "timer2")
	}()

	assert.True(t, <-timerEndChan)
	assert.True(t, timerEnd)
}
