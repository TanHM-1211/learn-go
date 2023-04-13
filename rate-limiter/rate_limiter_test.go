// go test -bench=. -benchmem -count=1
package ratelimiter

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func continuouslySpawn(num int, rl RateLimiter, timeInterval, sleepDuration time.Duration) {
	tick := time.Tick(timeInterval)
	allowed := 0
	denied := 0
	done := false
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		for range time.Tick(time.Second) {
			fmt.Printf("Allowed: %d; 	Denied: %d\n", allowed, denied)
			if done {
				wg.Done()
				break
			}
		}
	}()

	for i := 0; i < num; i++ {
		<-tick
		if rl.allow(i) {
			wg.Add(1)
			go func() {
				time.Sleep(sleepDuration)
				rl.done(i)
				wg.Done()
			}()
			allowed++
		} else {
			denied++
		}
		// if i%1000 == 0 {
		// 	fmt.Printf("Iter: %d;	Allowed: %d; 	Denied: %d\n", i, allowed, denied)
		// }
	}
	done = true
	wg.Wait()
}

func TestBasicRateLimiter(t *testing.T) {
	rl := NewBasicRateLimiter(time.Second, 1000)
	continuouslySpawn(10000, rl, 500*time.Microsecond, 900*time.Millisecond)

}
