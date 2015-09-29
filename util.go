package main

import (
	"math"
)

// pulled from https://gist.github.com/DavidVaini/10308388#gistcomment-1391788
func round(f float64) int {
	return int(math.Floor(f + .5))
}
