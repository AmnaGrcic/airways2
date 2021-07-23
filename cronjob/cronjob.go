package cronjob

import (
	"fmt"

	"github.com/jasonlvhit/gocron"
)

func CronJobs() {
	fmt.Println("cronjobs")

	gocron.Every(1).Minute().Do(Print)
	gocron.Every(10).Seconds().Do(Print)

	// gocron.Every(1).Second().Do(PrintWithParams, 1, "Hello")

	// gocron.Every(1).Monday().Do(Print)

	<-gocron.Start()
}

func Test() {
	fmt.Println("Test")
}

func Print() {
	fmt.Println("I am running task.")
}

// func PrintWithParams(a int, b string) {
// 	fmt.Println(a, b)
// }
