// Package rand is a helper package wrapping std lib functions.
package rand

import (
	"math/rand/v2"
	"time"
)

// Dur returns a random duration within the given min and max range.
func Dur(min, max time.Duration) time.Duration {
	if min >= max {
		return min
	}
	delta := max - min
	return min + time.Duration(rand.Int64N(int64(delta)))
}

// Bool returns a random boolean value.
func Bool() bool { return rand.IntN(2) == 1 }
