//https://go.dev/play/p/1CDa2j4GrV3
/*
Fan Out Pattern
The main idea behind Fan Out Pattern is to have:

a channel that provides a signaling semantics
channel can be buffered, so we don't wait on immediate receive confirmation
a goroutine that starts multiple (other) goroutines to do some work
a multiple goroutines that do some work and use signaling channel to signal that the work is done
*/

//In this example we have multiple employees that have some work to do. We also have a manager (main goroutine) that waits on that work to be done. Once each employee work is done, employee notifies manager by sending a signal (paper) via communication channel ch.

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	employees := 3

	// make channel of type string which provides signaling semantics
	// buffered channel is used so no goroutine blocks a sending operation
	// if two goroutines send a signal at the same time, channel performs synchronization
	ch := make(chan string, employees)

	for e := 0; e < employees; e++ {

		// start goroutine that does some work for employee e
		go func(employee int) {

			// simulate the idea of unknown latency (do not use in production)
			time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
			// when work is done send signal
			ch <- "paper"
			fmt.Println("employee : sent signal :", employee)
		}(e)
	}

	// goroutine 'main' => manager
	// goroutines 'main' and employee goroutines are executed in parallel

	// wait for all employee work to be done
	for employees > 0 {
		// receive signal sent from the employee
		p := <-ch

		employees--
		fmt.Println(p)
		fmt.Println("manager : received signal :", p)
	}

	time.Sleep(time.Second)
}

/*
Result
go run main.go

employee : sent signal : 1
paper
manager : received signal : paper
employee : sent signal : 2
paper
manager : received signal : paper
employee : sent signal : 0
paper
manager : received signal : paper
*/
