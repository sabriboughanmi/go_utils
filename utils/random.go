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
//	- noRewardsPossibility : True and percentages sum is less than 100% it sets it to 100%.
//	which means it gives the possibility to return -1 index. (No Element Selected)
func PeekRandomElement(percentages []float32, noRewardsPossibility bool) int {

	var totalPoints float32 = 0

	//Sum percentages to make a random choice
	for _, percentage := range percentages {
		totalPoints += percentage
	}

	// No specified rarity card is available to the user
	if totalPoints == 0 {
		return 0
	}

	var randomIndex = 0
	//Make no selected element possible
	if totalPoints < 100 && noRewardsPossibility {
		totalPoints = 100
		randomIndex = -1
	}

	//ensure that rand has a new seed value
	time.Sleep(time.Microsecond)
	randomRewardValue := RandFloat32(0, totalPoints, time.Now().UnixNano())

	// Pick a random Card
	var accumulatedPoints float32 = 0

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
