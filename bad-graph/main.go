package main

// A single function that sums the contents of an array, to explore generating and optimizing graphs
func main() {
	var array [8]int
	sum := 0
	for i := 0; i < 8; i++ {
		sum = array[i] + sum
	}
}
