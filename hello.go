package main

import (
	"fmt"
)

func main() {
	var seconds, hours, minutes int
	fmt.Scan(&seconds)
	hours = seconds / 30
	minutes = 2 * (seconds % 30)
	fmt.Println("It is", hours, "hours", minutes, "minutes.")
}
