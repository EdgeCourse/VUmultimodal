/*
Go closure middleware

The middleware are functions that execute during the lifecycle of a request to a server. The middleware is commonly used for logging, error handling, etc

In Go, middleware is often created with the help of closures.


*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	//handlefunc: where to go, what to do
	http.HandleFunc("/now", logDuration(getTime))
	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func logDuration(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		f(w, r)
		end := time.Now()
		fmt.Println("The request took", end.Sub(start))
	}
}

func getTime(w http.ResponseWriter, r *http.Request) {

	now := time.Now()
	_, err := fmt.Fprintf(w, "%s", now)

	if err != nil {
		log.Fatal(err)
	}
}
