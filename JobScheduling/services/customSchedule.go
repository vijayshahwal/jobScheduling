package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/vijayshahwal/jobScheduling/interfaces"
	"github.com/vijayshahwal/jobScheduling/models"
)

type CustomScheduleProcessor struct {
	cacheRepo interfaces.CacheRepository
}

func NewCustomScheduleProcessor(cacheRepo interfaces.CacheRepository) interfaces.ScheduleProcessor {
	return &CustomScheduleProcessor{cacheRepo: cacheRepo}
}

func (csp *CustomScheduleProcessor) CanProcess(schedule interface{}) bool {
	if dataStr, ok := schedule.(string); ok {
		var rawData map[string]interface{}
		if err := json.Unmarshal([]byte(dataStr), &rawData); err != nil {
			return false
		}

		return csp.isCustomSchedule(rawData)
	}

	_, ok := schedule.(models.CustomSchedule)
	return ok
}
func (csp *CustomScheduleProcessor) Process(schedule interface{}, jobKey string) error {
	var customSchedule models.CustomSchedule

	if dataStr, ok := schedule.(string); ok {
		if err := json.Unmarshal([]byte(dataStr), &customSchedule); err != nil {
			return err
		}
	} else {
		var ok bool
		customSchedule, ok = schedule.(models.CustomSchedule)
		if !ok {
			return fmt.Errorf("invalid schedule type")
		}
	}

	currentTime := time.Now().Unix()

	if customSchedule.NextInvocation == 0 || currentTime >= customSchedule.NextInvocation {
		fmt.Println("Custom Job", customSchedule.JobId, " - Scheduled and Completed")

		customSchedule.NextInvocation = customSchedule.CalculateNextRun()
		return csp.cacheRepo.Set(context.Background(), jobKey, customSchedule, 0)
	}

	return nil
}

func (csp *CustomScheduleProcessor) isCustomSchedule(rawData map[string]interface{}) bool {
	// CustomSchedule-specific fields that must be present
	customSpecificFields := []string{"dayOfMonth", "dayOfWeek", "month", "year"}

	// Check that at least some CustomSchedule-specific fields are present
	hasCustomFields := false
	for _, field := range customSpecificFields {
		if _, exists := rawData[field]; exists {
			hasCustomFields = true
			break
		}
	}

	if !hasCustomFields {
		return false // This is not a CustomSchedule
	}

	// Check that jobId is present
	if _, exists := rawData["jobId"]; !exists {
		return false
	}

	// Make sure this is not a FixedSchedule by checking for "daily" field
	if _, exists := rawData["daily"]; exists {
		return false // This has "daily" field, so it's a FixedSchedule
	}

	return true
}
