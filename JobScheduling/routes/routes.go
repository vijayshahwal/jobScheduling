package routes

import (
	"github.com/vijayshahwal/jobScheduling/controllers"

	"github.com/gorilla/mux"
)

func RegisterJobRoutes(router *mux.Router, jobController *controllers.JobController) {
	// Job management routes
	router.HandleFunc("/job", jobController.CreateJob).Methods("POST")
	router.HandleFunc("/job", jobController.GetAllJobs).Methods("GET")
	router.HandleFunc("/job/{jobId}", jobController.GetJobById).Methods("GET")

	// Job scheduling routes
	router.HandleFunc("/job/{jobId}/schedule/fixed", jobController.ScheduleFixedJob).Methods("POST")
	router.HandleFunc("/job/{jobId}/schedule/custom", jobController.ScheduleCustomJob).Methods("POST")
}
