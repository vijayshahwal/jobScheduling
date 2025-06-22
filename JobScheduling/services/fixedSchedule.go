package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/vijayshahwal/jobScheduling/interfaces"
	"github.com/vijayshahwal/jobScheduling/models"
)

type FixedScheduleProcessor struct {
	cacheRepo interfaces.CacheRepository
}

func NewFixedScheduleProcessor(cacheRepo interfaces.CacheRepository) interfaces.ScheduleProcessor {
	return &FixedScheduleProcessor{cacheRepo: cacheRepo}
}

func (fsp *FixedScheduleProcessor) CanProcess(schedule interface{}) bool {
	if dataStr, ok := schedule.(string); ok {
		var rawData map[string]interface{}
		if err := json.Unmarshal([]byte(dataStr), &rawData); err != nil {
			return false
		}

		// Check if this is specifically a FixedSchedule by validating required fields
		// and ensuring CustomSchedule-specific fields are NOT present
		return fsp.isFixedSchedule(rawData)
	}

	_, ok := schedule.(models.FixedSchedule)
	return ok
}

// Helper function to validate FixedSchedule struct
func (fsp *FixedScheduleProcessor) validateFixedScheduleStruct(schedule models.FixedSchedule) error {
	if schedule.JobId == "" {
		return fmt.Errorf("jobId cannot be empty")
	}

	if schedule.Minutes == "" && schedule.Hours == "" && schedule.Daily == "" {
		return fmt.Errorf("at least one of 'minutes', 'hours', or 'daily' must be greater than 0")
	}

	return nil
}

func (fsp *FixedScheduleProcessor) Process(schedule interface{}, jobKey string) error {
	var fixedSchedule models.FixedSchedule
	var err error

	if dataStr, ok := schedule.(string); ok {
		// fmt.Println(string(dataStr))
		err = json.Unmarshal([]byte(dataStr), &fixedSchedule)
		if err != nil {
			return fmt.Errorf("failed to unmarshal FixedSchedule: %w", err)
		}

		// Validate the unmarshaled struct
		if err = fsp.validateFixedScheduleStruct(fixedSchedule); err != nil {
			return fmt.Errorf("invalid FixedSchedule data: %w", err)
		}
	} else if fs, ok := schedule.(models.FixedSchedule); ok {
		if err = fsp.validateFixedScheduleStruct(fs); err != nil {
			return fmt.Errorf("invalid FixedSchedule struct: %w", err)
		}
		fixedSchedule = fs
	} else {
		return fmt.Errorf("invalid schedule type, expected FixedSchedule")
	}

	currentTime := time.Now().Unix()

	if fixedSchedule.NextInvocation == 0 || currentTime >= fixedSchedule.NextInvocation {
		fmt.Printf("Fixed Job %s - Scheduled and Completed - is scheduled for every Minutes: %s, Hours: %s, Daily: %s\n",
			fixedSchedule.JobId, fixedSchedule.Minutes, fixedSchedule.Hours, fixedSchedule.Daily)

		fixedSchedule.NextInvocation = fixedSchedule.CalculateNextRun()
		return fsp.cacheRepo.Set(context.Background(), jobKey, fixedSchedule, 0)
	}

	return nil
}

func (fsp *FixedScheduleProcessor) isFixedSchedule(rawData map[string]interface{}) bool {
	// Required fields for FixedSchedule
	requiredFields := []string{"jobId", "minutes", "hours", "daily"}

	// Fields that indicate this is a CustomSchedule (should NOT be present)
	customScheduleFields := []string{"dayOfMonth", "dayOfWeek", "month", "year"}

	// Check if any CustomSchedule-specific fields are present
	for _, field := range customScheduleFields {
		if _, exists := rawData[field]; exists {
			return false // This is a CustomSchedule, not a FixedSchedule
		}
	}

	// Check that all required FixedSchedule fields are present
	for _, field := range requiredFields {
		if _, exists := rawData[field]; !exists {
			return false
		}
	}

	// Validate that daily field exists (specific to FixedSchedule)
	if _, exists := rawData["daily"]; !exists {
		return false
	}

	return true
}
