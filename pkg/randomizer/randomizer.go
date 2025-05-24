package randomizer

import (
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func Randomize[T any](enums ...T) T {
	if len(enums) == 0 {
		var zeroValue T
		return zeroValue
	}

	return enums[rng.Intn(len(enums)-1)]
}
