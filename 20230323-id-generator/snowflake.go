package main

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
)

const (
	BitLenTime          = 39
	BitLenMachineID     = 10
	BitLenSequence      = 63 - BitLenTime - BitLenMachineID
	TimeWindowInMiliSec = 10
	BitShiftMachineID   = BitLenSequence
	BitShiftTime        = BitShiftMachineID + BitLenMachineID
)

var (
	MaxTime     int64 = int64(math.Pow(float64(2), float64(BitLenTime))) - 1
	MaxSequence int64 = int64(math.Pow(float64(2), float64(BitLenSequence))) - 1

	// epoch is the start time
	epoch int64 = time.Date(2023, 3, 23, 0, 0, 0, 0, time.Now().Location()).UnixMilli() / 10
)

var (
	ErrTimeExceeded     = errors.New("Time Exceed Error")
	ErrSequenceExceeded = errors.New("Sequence Exceed Error")
)

type ID int64

type Generator struct {
	mu                       *sync.Mutex
	startTime                int64
	elapsedTime              int64
	sequence                 int64
	machineIDsequenceShifted int64
}

func newGenerator(startTime int64, machineID uint16) *Generator {
	return &Generator{
		mu:                       new(sync.Mutex),
		startTime:                startTime,
		elapsedTime:              time.Now().UnixMilli() / 10,
		sequence:                 0,
		machineIDsequenceShifted: int64(machineID) << BitShiftMachineID,
	}
}

func createID(currentTime int64, machineIDsequenceShifted int64, sequence int64) (ID, error) {
	if currentTime > MaxTime {
		return 0, ErrTimeExceeded
	} else if sequence > MaxSequence {
		return 0, ErrSequenceExceeded
	} else {
		id := ID((currentTime << BitShiftTime) | machineIDsequenceShifted | sequence)
		fmt.Println(id, currentTime, machineIDsequenceShifted, sequence)
		return id, nil
	}
}

func (g *Generator) generate() (ID, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	var currentTime int64 = time.Now().UnixMilli() / 10
	if currentTime == g.elapsedTime {
		g.sequence += 1
	} else {
		g.sequence = 0
		g.elapsedTime = currentTime
	}

	id, err := createID(g.elapsedTime-epoch, g.machineIDsequenceShifted, g.sequence)
	if err != nil {
		return 0, err
	} else {
		return id, nil
	}
}

func (id ID) getInfo() {
	timeStamp := id >> BitShiftTime
	machineID := (id >> BitShiftMachineID) & ((1 << BitLenMachineID) - 1)
	sequence := id & ((1 << BitLenSequence) - 1)
	fmt.Printf("id: %d,\t time: %d,\t machineID: %d,\t sequence:%d\n", id, timeStamp, machineID, sequence)
}

func testSequentialRequest(g *Generator, numTest int) {
	result := make(map[ID]bool)
	var id ID
	var err error
	for i := 0; i < numTest; i++ {
		id, err = g.generate()
		if err != nil {
			panic(err)
		} else {
			if _, ok := result[id]; ok {
				panic(fmt.Sprintf("%d already existed", id))
			} else {
				result[id] = true
			}
		}
	}
	for k, _ := range result {
		k.getInfo()
	}
}

func main() {
	generator := newGenerator(time.Now().Unix(), 1)
	testSequentialRequest(generator, 10000)
}
