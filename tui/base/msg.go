package base

import "time"

type BackMsg struct{}

type NotificationMsg struct {
	Message string
}

type NotificationWithDurationMsg struct {
	NotificationMsg
	Duration time.Duration
}

type NotificationTimeoutMsg struct{}
