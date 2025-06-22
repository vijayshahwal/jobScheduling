package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/vijayshahwal/jobScheduling/container"
	"github.com/vijayshahwal/jobScheduling/interfaces"
	"github.com/vijayshahwal/jobScheduling/routes"
)

func main() {
	container := container.NewContainer()

	r := mux.NewRouter()
	routes.RegisterJobRoutes(r, container.JobController)

	// Start scheduler in background
	go StartScheduler(container.ScheduleService)

	fmt.Println("Server started listening on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}

func StartScheduler(scheduleService interfaces.ScheduleService) {
	for {
		if err := scheduleService.ProcessSchedules(context.Background()); err != nil {
			log.Printf("Error processing schedules: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
}
