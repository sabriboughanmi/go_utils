package utils

import (
	"math/rand"
	"time"
)

//GetRandomNumberInRange returns a random number n between min <= n < max.
//this function sleeps for a nanosecond before getting the random number in order to guarantee the number authenticity. (Slow)
func GetRandomNumberInRange(min, max int) int {
	time.Sleep(1 * time.Nanosecond)
	return GetRandomNumberInRangeWithSeed(min, max,time.Now().UnixNano() )
}

//GetRandomNumberInRangeWithSeed returns a random number n between min <= n < max using a seed value.
func GetRandomNumberInRangeWithSeed(min, max int, seed int64) int {
	rand.Seed(seed)
	return rand.Intn(max-min+1) + min
}

