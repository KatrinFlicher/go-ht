package main

func Filter(slice []int, fn func(item, index int) bool) (result []int) {
	for i, value := range slice {
		if fn(value, i) {
			result = append(result, value)
		}
	}
	return
}

func main() {
}
