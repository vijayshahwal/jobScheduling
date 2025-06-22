package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vijayshahwal/jobScheduling/interfaces"
	"github.com/vijayshahwal/jobScheduling/models"
)

type JobController struct {
	jobService      interfaces.JobService
	scheduleService interfaces.ScheduleService
}

func NewJobController(jobService interfaces.JobService, scheduleService interfaces.ScheduleService) *JobController {
	return &JobController{
		jobService:      jobService,
		scheduleService: scheduleService,
	}
}

func (jc *JobController) CreateJob(w http.ResponseWriter, r *http.Request) {
	var job models.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	createdJob, err := jc.jobService.CreateJob(context.Background(), job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdJob)
}

func (jc *JobController) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := jc.jobService.GetAllJobs(context.Background())
	if err != nil {
		http.Error(w, "Error fetching jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func (jc *JobController) GetJobById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobId := vars["jobId"]

	job, err := jc.jobService.GetJob(context.Background(), jobId)
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (jc *JobController) ScheduleFixedJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobId := vars["jobId"]

	var schedule models.FixedSchedule
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := jc.scheduleService.ScheduleFixedJob(context.Background(), jobId, schedule); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Job scheduled successfully"})
}

func (jc *JobController) ScheduleCustomJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobId := vars["jobId"]

	var schedule models.CustomSchedule
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := jc.scheduleService.ScheduleCustomJob(context.Background(), jobId, schedule); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Job scheduled successfully"})
}
