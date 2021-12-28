package main

import (
	"fmt"
)

func main() {
	var number int
	fmt.Scan(&number)
	switch {
	case number == 10000:
		fmt.Println(1)
	case 10000 > number && number > 999:
		fmt.Println(number / 1000)
	case 1000 > number && number > 99:
		fmt.Println(number / 100)
	case 100 > number && number > 9:
		fmt.Println(number / 10)
	default:
		fmt.Println(number)
	}
}
