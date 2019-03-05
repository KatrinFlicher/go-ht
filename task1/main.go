package main

import "fmt"

func Filter(slice []int, fn func(item, index int) bool) (result []int) {
	for i := 0; i < len(slice); i++ {
		if fn(slice[i], i) {
			result = append(result, slice[i])
		}
	}
	return
}

func main() {
	fmt.Println("Even", Filter([]int{1, 2, 3, 4, 5}, func(item, index int) bool { return item%2 == 0 }))
}
