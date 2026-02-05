package election

import "time"

type ElectionStatus string

const (
	Draft ElectionStatus = "draft"
	Active ElectionStatus = "active"
	Closed ElectionStatus = "closed"
	Archived ElectionStatus = "archived"
)

type Election struct {
	ID string
	Title string `binding:"required"`
	Description string `binding:"required"`
	StartTime time.Time `binding:"required"`
	EndTime time.Time `binding:"required"`
	Status ElectionStatus `binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (status ElectionStatus) IsValid() bool {
	switch status {
	case Draft, Active, Closed, Archived:
		return true
	}
	return false
}