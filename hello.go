package main

import (
	"fmt"
)

func main() {
	var number int
	fmt.Scan(&number)
	var first, second, third int = number / 100, (number / 10) % 10, number % 10
	if first != second && second != third && third != first {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
