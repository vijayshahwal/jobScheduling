package models

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron"
)

type CustomSchedule struct {
	Id             string `json:"id" bson:"_id,omitempty"`
	JobId          string `json:"jobId" bson:"jobId"`
	Minutes        string `json:"minutes" bson:"minutes"`
	Hours          string `json:"hours" bson:"hours"`
	DayOfMonth     string `json:"dayOfMonth" bson:"dayOfMonth"`
	DayOfWeek      string `json:"dayOfWeek" bson:"dayOfWeek"`
	Month          string `json:"month" bson:"month"`
	Year           string `json:"year" bson:"year"`
	NextInvocation int64  `json:"nextInvocation,omitempty" bson:"nextInvocation,omitempty"`
}

func (cs *CustomSchedule) CalculateNextRun() int64 {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	cronExpression := fmt.Sprintf("%s %s %s %s %s", cs.Minutes, cs.Hours, cs.DayOfMonth, cs.Month, cs.DayOfWeek)

	schedule, err := parser.Parse(cronExpression)
	if err != nil {
		log.Fatal("unable to parse cronExpression", cronExpression)
		return time.Now().Unix()
	}

	return schedule.Next(time.Now()).Unix()
}
