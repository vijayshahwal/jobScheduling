package interfaces

import (
	"context"

	"github.com/vijayshahwal/jobScheduling/models"
)

type JobService interface {
	CreateJob(ctx context.Context, job models.Job) (*models.Job, error)
	GetJob(ctx context.Context, id string) (*models.Job, error)
	GetAllJobs(ctx context.Context) ([]models.Job, error)
}

type ScheduleService interface {
	ScheduleFixedJob(ctx context.Context, jobID string, schedule models.FixedSchedule) error
	ScheduleCustomJob(ctx context.Context, jobID string, schedule models.CustomSchedule) error
	ProcessSchedules(ctx context.Context) error
}

type ValidationService interface {
	ValidateJob(job models.Job) error
	ValidateFixedSchedule(schedule models.FixedSchedule) error
	ValidateCustomSchedule(schedule models.CustomSchedule) error
}
