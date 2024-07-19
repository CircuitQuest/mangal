package download

import "time"

type nextChapterMsg struct{}

type retryChapterMsg struct {
	After time.Duration
}

type downloadCompletedMsg struct{}
