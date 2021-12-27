package main

import (
	"fmt"
)

func main() {
	var number int
	fmt.Scan(&number)
	switch {
	case number == 0:
		fmt.Println("Ноль")
	case number > 0:
		fmt.Println("Число положительное")
	default:
		fmt.Println("Число отрицательное")
	}
}
