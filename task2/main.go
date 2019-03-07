package main

func MapTo(slice []int, fn func(item, index int) string) (result []string) {
	for i, value := range slice {
		result = append(result, fn(value, i))
	}
	return
}

func Convert(arr []int) []string {
	return MapTo(arr, func(item, index int) (res string) {
		res = "unknown"
		var m = map[int]string{
			1: "one", 2: "two", 3: "three", 4: "four", 5: "five", 6: "six", 7: "seven", 8: "eight", 9: "nine"}
		value, ok := m[item]
		if ok {
			res = value
		}
		return
	})
}

func main() {

}
