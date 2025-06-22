package services

import (
	"context"
	"fmt"

	"github.com/vijayshahwal/jobScheduling/interfaces"
	"github.com/vijayshahwal/jobScheduling/models"
)

type ScheduleService struct {
	cacheRepo  interfaces.CacheRepository
	jobRepo    interfaces.JobRepository
	validator  interfaces.ValidationService
	processors []interfaces.ScheduleProcessor
}

func NewScheduleService(
	cacheRepo interfaces.CacheRepository,
	jobRepo interfaces.JobRepository,
	validator interfaces.ValidationService,
) interfaces.ScheduleService {
	return &ScheduleService{
		cacheRepo: cacheRepo,
		jobRepo:   jobRepo,
		validator: validator,
		processors: []interfaces.ScheduleProcessor{
			NewFixedScheduleProcessor(cacheRepo),
			NewCustomScheduleProcessor(cacheRepo),
		},
	}
}

func (ss *ScheduleService) ScheduleFixedJob(ctx context.Context, jobID string, schedule models.FixedSchedule) error {
	if err := ss.validator.ValidateFixedSchedule(schedule); err != nil {
		return err
	}

	// Verify job exists
	_, err := ss.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("job not found: %v", err)
	}

	schedule.JobId = jobID
	key := fmt.Sprintf("job:%s", jobID)
	return ss.cacheRepo.Set(ctx, key, schedule, 0)
}

func (ss *ScheduleService) ScheduleCustomJob(ctx context.Context, jobID string, schedule models.CustomSchedule) error {
	if err := ss.validator.ValidateCustomSchedule(schedule); err != nil {
		return err
	}

	// Verify job exists
	_, err := ss.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("job not found: %v", err)
	}

	schedule.JobId = jobID
	key := fmt.Sprintf("job:%s", jobID)
	return ss.cacheRepo.Set(ctx, key, schedule, 0)
}

func (ss *ScheduleService) ProcessSchedules(ctx context.Context) error {
	keys, err := ss.cacheRepo.Scan(ctx, "job:*")
	if err != nil {
		return err
	}

	for _, key := range keys {
		data, err := ss.cacheRepo.Get(ctx, key)
		if err != nil {
			continue
		}

		for _, processor := range ss.processors {
			processorType := fmt.Sprintf("%T", processor)
			canProcess := processor.CanProcess(data)

			if canProcess {
				if err := processor.Process(data, key); err != nil {
					fmt.Printf("Error processing schedule %s with %s: %v\n", key, processorType, err)
				}
				break
			}
		}
	}

	return nil
}
