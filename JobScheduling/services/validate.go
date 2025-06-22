package services

import (
	"fmt"

	"github.com/spf13/cast"
	"github.com/vijayshahwal/jobScheduling/interfaces"
	"github.com/vijayshahwal/jobScheduling/models"
)

type ValidationService struct{}

func NewValidationService() interfaces.ValidationService {
	return &ValidationService{}
}

func (vs *ValidationService) ValidateJob(job models.Job) error {
	if job.Name == "" {
		return fmt.Errorf("job name is required")
	}
	return nil
}

func (vs *ValidationService) ValidateFixedSchedule(schedule models.FixedSchedule) error {
	min := cast.ToInt(schedule.Minutes)
	hours := cast.ToInt(schedule.Hours)
	daily := cast.ToInt(schedule.Daily)

	if min == 0 && hours == 0 && daily == 0 {
		return fmt.Errorf("at least one time interval must be specified")
	}
	return nil
}

func (vs *ValidationService) ValidateCustomSchedule(schedule models.CustomSchedule) error {
	if schedule.Minutes == "" || schedule.Hours == "" {
		return fmt.Errorf("minutes and hours must be specified")
	}
	return nil
}
