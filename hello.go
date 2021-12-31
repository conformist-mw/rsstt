package main

import (
	"fmt"
)

func main() {
	var count, num, sum int
	fmt.Scan(&count)
	for i := 0; i < count; i++ {
		fmt.Scan(&num)
		if 10 <= num && num < 100 && num%8 == 0 {
			sum += num
		}
	}
	fmt.Println(sum)
}
