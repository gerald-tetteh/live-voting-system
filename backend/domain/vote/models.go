package vote

import "time"

type Vote struct {
	ID string
	ElectionId string
	UserId string
	CreatedAt time.Time
	UpdatedAt time.Time
}