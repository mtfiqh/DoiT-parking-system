package randomizer

import (
	"math/rand"
)

func RandomizeEnum[T any](enums ...T) T {

	if len(enums) == 0 {
		var zeroValue T
		return zeroValue
	}

	return enums[rand.Intn(len(enums))]
}

func RandomizeInt(min, max int) int {
	if min >= max {
		return min
	}
	return rand.Intn(max-min) + min + 1
}
