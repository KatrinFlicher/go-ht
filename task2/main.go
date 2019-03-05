package main

import "fmt"

func MapTo(slice []int, fn func(item, index int) string) (result []string) {
	for i := 0; i < len(slice); i++ {
		result = append(result, fn(slice[i], i))
	}
	return
}

func Convert(arr []int) []string {
	return MapTo(arr, func(item, index int) (res string) {
		switch item {
		case 1:
			res = "one"
		case 2:
			res = "two"
		case 3:
			res = "three"
		case 4:
			res = "four"
		case 5:
			res = "five"
		case 6:
			res = "six"
		case 7:
			res = "seven"
		case 8:
			res = "eight"
		case 9:
			res = "nine"
		default:
			res = "unknown"
		}
		return
	})
}

func main() {
	fmt.Println(MapTo([]int{1, 2, 3, 4, 5, 10}, func(item, index int) string { return item * 2 }))
}
