package go_countdown_timer

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMjTimerCanProcessAndStop(t *testing.T) {

	timer := NewTimer(3 * time.Second)

	countProcess := 0

	go func() {
		time.Sleep(time.Second * 1)
		timer.GetProcessChan() <- 0
	}()

	expired := false

	timer.StartWithActions(func(b interface{}) bool {
		countProcess++
		return false
	},
		func() {

			expired = true
		})
	assert.Equal(t, 1, countProcess)
	assert.True(t, expired)

}
func TestMjTimerCanProcessMultipleChannel(t *testing.T) {

	timer := NewTimer(3 * time.Second)

	marker := map[string]interface{}{
		"foo": nil,
		"bar": nil,
	}
	go func() {
		time.Sleep(time.Second * 2)
		timer.GetProcessChan() <- "foo"
		timer.GetProcessChan() <- "bar"
	}()
	timer.StartWithActions(func(arg interface{}) bool {

		id := arg.(string)
		if _, ok := marker[id]; !ok {
			return false
		}

		marker[id] = arg
		cur := 0
		for _, v := range marker {
			if v != nil {
				cur++
			}
		}

		if cur == len(marker) {
			fmt.Printf("collected done :%v \n", marker)

			return true
		}
		return false
	}, func() {
		fmt.Printf("process something after timeout")
	})
	fmt.Printf("marker: %v \n", marker)

	assert.Equal(t, "foo", marker["foo"])
	assert.Equal(t, "bar", marker["bar"])

}
func TestMjTimerCanProcessMultipleChannelAndTimeout(t *testing.T) {

	timer := NewTimer(3 * time.Second)

	marker := map[string]interface{}{
		"foo": nil,
		"bar": nil,
	}
	go func() {
		time.Sleep(time.Second * 2)
		timer.GetProcessChan() <- "foo"

	}()

	timer.StartWithActions(func(arg interface{}) bool {

		id := arg.(string)
		if _, ok := marker[id]; !ok {
			return false
		}

		marker[id] = arg
		cur := 0
		for _, v := range marker {
			if v != nil {
				cur++
			}
		}

		if cur == len(marker) {
			fmt.Printf("collected done :%v", marker)
			return true
		}
		return false
	}, func() {
		fmt.Printf("process something after timeout")
	})
	fmt.Printf("marker: %v", marker)
	assert.Equal(t, "foo", marker["foo"])
	assert.Equal(t, nil, marker["bar"])
}
