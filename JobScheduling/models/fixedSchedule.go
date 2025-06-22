package models

import (
	"time"

	"github.com/spf13/cast"
)

type FixedSchedule struct {
	Id             string `json:"id" bson:"_id,omitempty"`
	JobId          string `json:"jobId" bson:"jobId"`
	Minutes        string `json:"minutes" bson:"minutes"`
	Hours          string `json:"hours" bson:"hours"`
	Daily          string `json:"daily" bson:"daily"`
	Type           string `json:"type" bson:"type"`
	NextInvocation int64  `json:"nextInvocation,omitempty" bson:"nextInvocation,omitempty"`
}

func (fs *FixedSchedule) CalculateNextRun() int64 {
	currentTime := time.Now().Unix()

	if cast.ToInt(fs.Minutes) > 0 {
		return currentTime + cast.ToInt64(fs.Minutes)*60
	} else if cast.ToInt(fs.Hours) > 0 {
		return currentTime + cast.ToInt64(fs.Hours)*3600
	} else if cast.ToInt(fs.Daily) > 0 {
		return currentTime + cast.ToInt64(fs.Daily)*24*3600
	}

	return currentTime
}
