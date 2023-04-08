// go test -bench=. -benchmem -count=5
package BloomFilter

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789/+"
const (
	numBitsLetters    = 6
	letterIdxMask     = 1<<numBitsLetters - 1
	numLettersPerRand = 63 / numBitsLetters
)

func randString(n int) string {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := make([]byte, n)
	randNum := rand.Int63()
	remain := numBitsLetters
	for i := 0; i < n; i++ {
		res[i] = letterBytes[randNum&letterIdxMask]
		randNum >>= numBitsLetters
		remain -= 1
		if remain == 0 {
			randNum = rand.Int63()
			remain = numBitsLetters
		}
	}
	return string(res)
}

// func TestRandString(t *testing.T) {
// 	for i := 0; i < 10000; i++ {
// 		result := randString(i)
// 		if len(result) != i {
// 			t.Error("failed")
// 		}
// 	}
// }

// func BenchmarkRandString(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		randString(1000)
// 	}
// }

func runBloomFilter(t *testing.T, size int, errRate float32) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	bloomFilter := NewFilter(size, errRate)
	fmt.Println(size, errRate, len(bloomFilter.hashers), bloomFilter.size())
	count := 0
	var input string
	mark := make(map[string]bool)

	for i := 0; i < size; i++ {
		for {
			input = randString(rand.Intn(990) + 10)
			if _, ok := mark[input]; !ok {
				mark[input] = true
				break
			}
		}

		if bloomFilter.contains(input) {
			count += 1
		}
		bloomFilter.insert(input)
	}
	if float32(count)/float32(size) > 2*errRate {
		t.Errorf("error with size=%d and errRate=%f and count=%d", size, errRate, count)
	}
}

func TestBloomFilter(t *testing.T) {
	sizes := []int{1000, 10000, 100000}
	errRates := []float32{0.005, 0.01, 0.05, 0.1}
	for _, size := range sizes {
		for _, errRate := range errRates {
			runBloomFilter(t, size, errRate)
		}
	}
	runBloomFilter(t, 1000000, 0.01)
}
