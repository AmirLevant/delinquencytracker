package main

import (
	"fmt"
	"time"
)

func main() {
	currentTime := time.Now()
	var printy string = currentTime.Format(time.DateOnly)
	fmt.Println("the date is ", printy)
}
