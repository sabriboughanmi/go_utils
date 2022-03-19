package utils

import (
	"math/rand"
	"time"
)

func RandFloat32(min, max float32, seed int64) float32 {
	rand.Seed(seed)
	return min + rand.Float32()*(max-min)
}

//PeekRandomElement peeks a random element from a set of win rates percentages.
func PeekRandomElement(percentages []float32) int {

	var totalPoints float32 = 0

	//Sum percentages to make a random choice
	for _, percentage := range percentages {

		totalPoints += percentage
	}

	// No specified rarity card is available to the user
	if totalPoints == 0 {
		return 0
	}

	//ensure that rand has a new seed value
	time.Sleep(time.Microsecond)
	randomRewardValue := RandFloat32(0, totalPoints, time.Now().UnixNano())

	// Pick a random Card
	var accumulatedPoints float32 = 0

	var randomIndex = 0

	//Iterate over Cards
	for i, percentage := range percentages {
		accumulatedPoints += percentage
		if randomRewardValue < accumulatedPoints {
			randomIndex = i
			break
		}
	}

	return randomIndex
}
