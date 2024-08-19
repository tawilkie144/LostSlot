package Entities

import (
	"time"
)

type Event struct {
	Id             int64
	Title          string
	Description    string
	Location       string
	Time           TimeBlock
	Organizer      User
	Attendees      []User
	AttendeeStatus AttendeeStatus
	EventLocation  EventLocation
}

type EventLocation struct {
	LocationData   Location
	IsVirtual      bool
	VirtualAddress string
}

type TimeBlock struct {
	StartTime time.Time
	Duration  time.Duration
}

type AttendeeStatus int

const (
	Accepted AttendeeStatus = iota
	Declined
	Tentative
	NoResponse
)

func (status AttendeeStatus) String() string {
	return [...]string{"Accepted", "Declined", "Tentative", "No Response"}[status-1]
}

func (status AttendeeStatus) EnumId() int {
	return int(status)
}
