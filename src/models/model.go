package models

import (
	"time"
)

const (
	EVENT_BROAD = iota
	EVENT_UNICAST
)

type Event struct {
	Name      string
	Content   string
	Timestamp string
	SendType  int
}

func NewEvent(htype int, name, content string) Event {
	t := time.Now()
	return Event{
		Name:      name,
		Content:   content,
		Timestamp: t.Format(time.RFC3339),
		SendType:  htype,
	}
}
