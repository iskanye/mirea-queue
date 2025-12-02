package models

import "fmt"

type Queue struct {
	Group   string
	Subject string
}

func (q *Queue) Key() string {
	return fmt.Sprintf("%s : %s", q.Group, q.Subject)
}

type QueueEntry struct {
	ChatID string
}
