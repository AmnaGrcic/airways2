package main

import (
	"airways/repository"
	"airways/routes"
	"airways/redis"
)

func main() {

	repository.ConnectToDatabase()
	repository.Migration()
	// go cronjob.CronJobs()
	redis.StartRedis()
	routes.Startroutes()
}
