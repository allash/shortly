package generator

import (
	"math/rand/v2"
)

func rangeIn(low, hi int64) int64 {
	return low + rand.Int64N(hi - low)
}

func GenerateId() int64 {
	return rangeIn(10000000, 99999999)
}

