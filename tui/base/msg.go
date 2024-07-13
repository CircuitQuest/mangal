package base

import "time"

type BackMsg struct {
	Steps int
}

type BackToHomeMsg struct{}

type NotificationMsg struct {
	Message string
}

type NotificationWithDurationMsg struct {
	NotificationMsg
	Duration time.Duration
}

type NotificationTimeoutMsg struct{}

type LoadingMsg struct {
	Message string
}

// RestoredMsg is sent when going back to the state.
type RestoredMsg struct{}

type ShowLoadingMsg bool

type ShowSubtitleMsg bool
