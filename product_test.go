package demo

import (
	"fmt"
	"testing"
)

// Given an  array of integers, return an array of integers where the index is the product of all others in the array.
// example [2,4,5] would return [20,10,8], [4*5, 2*5, 2*4]
// example [1,3,5,4] => [1*3*5*4/1, 1*3*5*4/3, 1*3*5*4/5, 1*3*5*4/4]

func TestProducts(t *testing.T) {

	testData := []int{1, 2, 3, 4}

	outputData := make([]int, len(testData))

	for i := range testData {
		sum := 1
		for j := range testData {
			if i == j {
				continue
			}
			sum = sum * testData[j]
		}
		outputData[i] = sum
	}
	fmt.Println(outputData)
}
