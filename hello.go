package main

import (
	"fmt"
)

func main() {
	var (
		count   = 0
		max_num = 0
		num     = 1
	)
	for num != 0 {
		fmt.Scan(&num)
		if num > max_num {
			max_num = num
			count = 0
		}
		if num == max_num {
			count++
		}
	}
	fmt.Println(count)
}
