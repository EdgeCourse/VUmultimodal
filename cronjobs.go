//first is cron, second is gocron
//https://www.libhunt.com/compare-gocron-vs-cron

//Way one: cron v2

/*
cron.v2 is another popular package that allows you to implement scheduled jobs in Golang. Let’s implement the above example again, but this time with the help of the cron.v2 package.

Install the cron.v2 package by running the following command in our terminal:

go get gopkg.in/robfig/cron.v2

1. Imports the required packages
2. Creates a new scheduler instance s in the runCronJobs function using the cron.v2 package
3. Defines a cron job that runs every one second (@every 1s) and calls the hello function
4. Starts the scheduler in blocking mode, which blocks the current execution path
5. Calls the runCronJobs function inside the main function. However, since the runCronJobs function runs asynchronously, the execution falls through. The Scanln method, which requires you to press a key before the program exits, prevents this.

Running go run main.go, you should see a message printed to the terminal every second.

*/
package main

// 1
import (
	"fmt"

	"gopkg.in/robfig/cron.v2"
)

func hello(name string) {
	message := fmt.Sprintf("Hi, %v", name)
	fmt.Println(message)
}

func runCronJobs() {
	// 2
	s := cron.New()

	// 3
	s.AddFunc("@every 1s", func() {
		hello("Bob Loblaw")
	})

	// 4
	s.Start()
}

// 5
func main() {
	runCronJobs()
	fmt.Scanln()
}

//Way two: gocron
/*
Using the gocron Package

gocron is a job-scheduling package that lets you run Go functions at predetermined intervals by defining a simple, human-friendly syntax. Let’s start by writing a small Golang application that prints a message on the terminal in one-second intervals.

The code below:

1. Imports the required packages
2. Defines the hello function that prints a message and takes a string parameter called name
3. Creates a new scheduler instance s in the runCronJobs function using the gocron package
4. Defines a cron job that runs every one second (Every(1).Seconds()) and calls the hello function
5. Starts the scheduler in blocking mode (StartBlocking()), which blocks the current execution path
6. Calls the runCronJobs function inside the main function

Save the file and install the included packages (i.e., gocron) by running the following command in our terminal:

go mod tidy

This will add the gocron package to the dependencies in the go.mod and go.sum files.

go run main.go

You will see output every second.
*/

/*

package main

// 1
import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

// 2
func hello(name string) {
	message := fmt.Sprintf("Hi, %v", name)
	fmt.Println(message)
}

func runCronJobs() {
	// 3
	s := gocron.NewScheduler(time.UTC)

	// 4
	s.Every(1).Seconds().Do(func() {
		hello("Bob Loblaw")
	})

	// 5
	s.StartBlocking()
}

// 6
func main() {
	runCronJobs()
}
*/

/*
Limits of embedded schedulers in Golang:

While using packages like gocron and cron.v2 provides useful functionality, there are some limitations to be aware of.

First, you can’t improve what you can’t measure. Properly measuring stats like memory usage or total automation-execution time can help you speed up and optimize your workflows and these are not supported while using packages like gocron and cron.v2.

It's also essential to understand what's going on in the code. While this is less of an issue with small and simple workflows, cron makes it very difficult to track the data flow between different tasks as the complexity increases.

Additionally, in order to log the data, you have to write the logic for logging and find a way to persistently store your logs, which may require setting up a storage solution.

Failure handling and failure management are also crucial considerations. Cron runs strictly at specified times, meaning that if a task fails, it will only run again at the next scheduled time and there are no auto-retry mechanisms for failed jobs. This means that in order to use cron safely, you'll also need to build out some way to get notifications (e.g., via Slack or email) to receive alerts when schedules fail.

Finally, if your team is dependent on a code-based solution for writing automation, engineers may turn into bottlenecks and cause latency because nontechnical members of your team are unable to create and manage automation workflows without them.


An alternative is Airplane scheduled tasks:

Airplane is a developer-first SaaS platform with first-class support for running scheduled and nonscheduled operations. Using Airplane you can easily set up new jobs, debug scripts, manually execute tasks, and much more. While covering all the drawbacks of embedded schedulers in Golang, Airplane also provides audit logs, approval flows, failure notifications, and other features that promote security.


*/
