package Entities

import (
	"net/mail"
	"time"
)

type User struct {
	UserId               int64
	GivenName            string
	LastName             string
	Initials             string
	PreferredName        string
	HashedPassword       string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Email                mail.Address
	SubscribedEvents     []Event
	NotificationSettings PreferredNotificationMethod
	PrimaryAddress       Location
}

type Location struct {
	Address string
	Name    string
}

type PreferredNotificationMethod int

const (
	//Enum names here
	NoNotification PreferredNotificationMethod = 0
	Email          PreferredNotificationMethod = 1
	SMS            PreferredNotificationMethod = 2
)

func (T PreferredNotificationMethod) ToString() string {
	return [...]string{"Email", "SMS", "No Notification"}[T-1]
}

func (T PreferredNotificationMethod) GetIdNum() int {
	return int(T)
}
