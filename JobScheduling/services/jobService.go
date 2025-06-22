package services

import (
	"context"

	"github.com/vijayshahwal/jobScheduling/interfaces"
	"github.com/vijayshahwal/jobScheduling/models"
)

type JobService struct {
	jobRepo   interfaces.JobRepository
	validator interfaces.ValidationService
}

func NewJobService(jobRepo interfaces.JobRepository, validator interfaces.ValidationService) interfaces.JobService {
	return &JobService{
		jobRepo:   jobRepo,
		validator: validator,
	}
}

func (js *JobService) CreateJob(ctx context.Context, job models.Job) (*models.Job, error) {
	if err := js.validator.ValidateJob(job); err != nil {
		return nil, err
	}
	return js.jobRepo.Save(ctx, job)
}

func (js *JobService) GetJob(ctx context.Context, id string) (*models.Job, error) {
	return js.jobRepo.FindByID(ctx, id)
}

func (js *JobService) GetAllJobs(ctx context.Context) ([]models.Job, error) {
	return js.jobRepo.FindAll(ctx)
}
