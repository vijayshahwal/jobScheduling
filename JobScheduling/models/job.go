package models

import (
	"time"
)

type Job struct {
	Id          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Priority    int       `json:"priority" bson:"priority"`
	Type        string    `json:"type" bson:"type"`
	CreatedOn   time.Time `json:"createdOn" bson:"createdOn"`
}
