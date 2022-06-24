/*
const (
    C1 = iota + 1
    _
    C3
    C4
)
fmt.Println(C1, C3, C4) // "1 3 4"
*/
/*
func SumProduct(a, b int) (int, int) {
    return a+b, a*b
}

func main() {
    // I only want the sum, but not the product
    sum, _ := SumProduct(1,2) // the product gets discarded
    fmt.Println(sum) // prints 3
}
*/
/*
   pets := []string{"dog", "cat", "fish"}

   // Range returns both the current index and value
   // but sometimes you may only want to use the value
   for _, pet := range pets {
       fmt.Println(pet)
   }
*/

//in half-written program
//Silence the compiler
//It can be used to during development to avoid compiler errors about unused imports and variables in a half-written program.

package main

import (
	"log"
	"os"
	"fmt"
)

var _ = fmt.Printf // DEBUG: delete when done

func main() {
	f, err := os.Open("gocon.go")
	if err != nil {
		log.Fatal(err)
	}
	_ = f // TODO: read file
}
