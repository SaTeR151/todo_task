package entity

import "time"

// FromColumnID.Name -> ToColumnID.Name [When: Timestamp] [What: TaskID.Name]
type MoveEvent struct {
	ID           string
	TaskID       string
	FromColumnID string
	ToColumnID   string
	Timestamp    time.Time
}

type MoveEvents []MoveEvent

type MoveEventCreate struct {
	TaskID       string
	FromColumnID string
	ToColumnID   string
}

type GetMoveEventsOpts struct {
	TaskID string
}
