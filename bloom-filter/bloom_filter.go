package BloomFilter

import (
	"hash/fnv"
	"math"
	"math/rand"
	"time"
)

type Hasher struct {
	seed uint32
}

func (hs Hasher) hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32() + hs.seed
}

type BloomFilter struct {
	maxItems int
	data     []bool
	hashers  []*Hasher
}

func (bf *BloomFilter) size() int {
	return len(bf.data)
}

func (bf *BloomFilter) insert(s string) {
	for _, hasher := range bf.hashers {
		bf.data[int(hasher.hash(s))%bf.size()] = true
	}
}

func (bf *BloomFilter) contains(s string) bool {
	for _, hasher := range bf.hashers {
		if !bf.data[int(hasher.hash(s))%bf.size()] {
			return false
		}
	}
	return true
}

func NewFilter(maxItems int, fpRate float32) *BloomFilter {
	// https://en.wikipedia.org/wiki/Bloom_filter#Probability_of_false_positives
	// SEED := 2
	SEED := time.Now().UnixNano()
	NUM_HASHES := int(-math.Floor(math.Log2(float64(fpRate))))
	r := rand.New(rand.NewSource(int64(SEED)))

	hashers := []*Hasher{}
	for i := 0; i < NUM_HASHES; i++ {
		hashers = append(hashers, &Hasher{r.Uint32()})
	}

	bloomSize := int(-math.Floor(1.44 * math.Log2(float64(fpRate)) * float64(maxItems)))
	bloomData := make([]bool, bloomSize)
	return &BloomFilter{maxItems: maxItems, data: bloomData, hashers: hashers}
}
