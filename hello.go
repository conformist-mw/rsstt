package main

import (
	"fmt"
)

func main() {
	var number int
	for {
		fmt.Scan(&number)
		if number > 100 {
			break
		} else if number < 10 {
			continue
		} else {
			fmt.Println(number)
		}
	}
}
