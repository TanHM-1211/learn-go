package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func (id ID) getInfo() {
	timeStamp := id >> BitShiftTime
	machineID := (id >> BitShiftMachineID) & ((1 << BitLenMachineID) - 1)
	sequence := id & ((1 << BitLenSequence) - 1)
	fmt.Printf("id: %d,\t time: %d,\t machineID: %d,\t sequence:%d\n", id, timeStamp, machineID, sequence)
}

func TestSequentialRequest(t *testing.T) {
	g := newGenerator(time.Now().Unix(), 1)
	numTest := 10000
	result := make(map[ID]bool)
	var id ID
	var err error
	for i := 0; i < numTest; i++ {
		id, err = g.generate()
		if err != nil {
			t.Error(err)
		} else {
			if _, ok := result[id]; ok {
				t.Errorf("%d already existed", id)
			} else {
				result[id] = true
			}
		}
	}

	// for k, _ := range result {
	// 	k.getInfo()
	// }
}

type ConcurrentResult struct {
	mu  *sync.Mutex
	ids map[ID]bool
}

func request(g *Generator, result *ConcurrentResult, t *testing.T) {
	id, err := g.generate()
	if err != nil {
		t.Error(err)
	} else {
		result.mu.Lock()
		defer result.mu.Unlock()
		if _, ok := result.ids[id]; ok {
			t.Errorf("%d already existed", id)
		} else {
			result.ids[id] = true
		}
	}
}

func TestConcurrentRequest(t *testing.T) {
	g := newGenerator(time.Now().Unix(), 1)
	numGoroutines := 400
	numRequests := 1000000
	result := &ConcurrentResult{
		mu:  new(sync.Mutex),
		ids: make(map[ID]bool),
	}
	wg := new(sync.WaitGroup)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(result *ConcurrentResult, t *testing.T, gid int) {

			fmt.Printf("Goroutine number %d \n", gid)
			for j := 0; j < numRequests/numGoroutines; j++ {
				request(g, result, t)
				time.Sleep(time.Millisecond)
			}
			wg.Done()
		}(result, t, i)
	}
	wg.Wait()
}

// GOMAXPROCS=4 go test *.go -v
