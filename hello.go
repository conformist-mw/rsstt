package main

import (
	"fmt"
	"math"
)

func main() {
	var number int
	fmt.Scan(&number)
	var (
		first int = number % int(math.Pow10(6)) / int(math.Pow10(5))
		second int = number % int(math.Pow10(5)) / int(math.Pow10(4))
		third int = number % int(math.Pow10(4)) / int(math.Pow10(3))
		fourth int = number % int(math.Pow10(3)) / int(math.Pow10(2))
		fifth int = number % int(math.Pow10(2)) / int(math.Pow10(1))
		sixth int = number % int(math.Pow10(1)) / int(math.Pow10(0))
	)
	if (first + second + third) == (fourth + fifth + sixth) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
