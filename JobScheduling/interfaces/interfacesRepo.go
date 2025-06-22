package interfaces

import (
	"context"
	"time"

	"github.com/vijayshahwal/jobScheduling/models"
)

type JobRepository interface {
	Save(ctx context.Context, job models.Job) (*models.Job, error)
	FindByID(ctx context.Context, id string) (*models.Job, error)
	FindAll(ctx context.Context) ([]models.Job, error)
}

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	Scan(ctx context.Context, pattern string) ([]string, error)
	Delete(ctx context.Context, key string) error
}

type ScheduleRepository interface {
	SaveSchedule(ctx context.Context, jobID string, schedule interface{}) error
	GetSchedule(ctx context.Context, jobID string) (interface{}, error)
	GetAllSchedules(ctx context.Context) ([]interface{}, []string, error)
}
