package main

import (
	"fmt"
)

//var _ = fmt.Printf // DEBUG: delete when done

func main() {
	var timeZone = map[string]int{
		"UTC": 0 * 60 * 60,
		"EST": -5 * 60 * 60,
		"CST": -6 * 60 * 60,
		"MST": -7 * 60 * 60,
		"PST": -8 * 60 * 60,
	}
	/*
		var val1 = timeZone["PST"]
		var val2 = timeZone["invalid"]
		fmt.Printf("val1: %d val2: %d\n", val1, val2)
	*/
/*
	var val1, ok1 = timeZone["PST"]
	var val2, ok2 = timeZone["invalid"]
	fmt.Printf("val1: %d ok1: %t val2: %d ok2: %t\n", val1, ok1, val2, ok2)
*/
	
		if seconds, ok := timeZone["PST"]; ok {
			fmt.Println(seconds)
		}
	

}
