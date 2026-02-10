package main

import (
	"fmt"
)

func MovingAverage(data []float64, windowSize int) []float64 {
	if windowSize <= 0 || windowSize > len(data) {
		return nil
	}

	result := make([]float64, len(data)-windowSize+1)
	for i := 0; i < len(result); i++ {
		sum := 0.0
		for j := 0; j < windowSize; j++ {
			sum += data[i+j]
		}
		result[i] = sum / float64(windowSize)
	}
	return result
}

func main() {
	sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0}
	averages := MovingAverage(sampleData, 3)
	fmt.Println("Moving averages:", averages)
}