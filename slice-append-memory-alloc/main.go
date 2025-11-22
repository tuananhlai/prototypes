package main

import (
	"fmt"
	"math"
)

func main() {
	var data []int

	var lastCap int
	for i := range 100000 {
		data = append(data, i)
		newCap := cap(data)

		if lastCap != newCap {
			fmt.Printf("index = %d, newCap = %d, percentChange = %v%%\n", i, newCap, getPercentChange(lastCap, newCap))
			lastCap = newCap
		}
	}
}

func getPercentChange(lastCap, newCap int) float64 {
	if lastCap == 0 {
		return math.Inf(1)
	}

	percentChange := (float64(newCap-lastCap) * 100) / float64(lastCap)
	return math.Round(percentChange*100) / 100
}
