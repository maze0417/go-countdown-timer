# go-countdown-timer
Count down timer written by Go . It could collect actions before timeup


# How to Use


First ,Just create a new  countdown timer
```
timer := NewTimer(3 * time.Second)
```

Start counting down, and inject a timeout func , that will execute when time up
```
timer.StartNew(func() {
        //do stuff when timeout		
	})

```

Send event to timer, that will execute inside timer 
```
timer.ReceiveProcessEvent() <- func() bool {
            //true will end timer immediately
            //false will continue to collect event
			return true
		}
```
# Installation
``
go get -u github.com/maze0417/go-countdown-timer
``
# Notice
If we call StartNew(..) twice at the different go-routine , the first launch will be forced end up   

