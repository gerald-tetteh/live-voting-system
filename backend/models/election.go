package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Election struct {
	ID string
	Name string `binding:"required"`
	Winner string
	StartTime time.Time `binding:"required"`
	EndTime time.Time `binding:"required"`
	Completed bool
	createdAt time.Time
	updatedAt time.Time
}

func New(name string, startTime time.Time, endTime time.Time) (*Election, error) {
	if startTime.Before(time.Now()) {
		return nil, errors.New("'startTime' must be in the future")
	}
	if name == "" || endTime.Before(startTime) || startTime.Equal(endTime) {
		return nil, errors.New("Election must have a name and the 'endTime' must be after 'startTime'")
	}
	currentTime := time.Now()
	return &Election{
		ID: uuid.NewString(),
		Name: name,
		StartTime: startTime,
		EndTime: endTime,
		Completed: false,
		createdAt: currentTime,
		updatedAt: currentTime,
	}, nil
}

func (election *Election) Save() {}