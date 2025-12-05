package models

import (
	"fmt"
	"strings"
)

type Queue struct {
	Group   string
	Subject string
}

func QueueFromKey(key string) Queue {
	data := strings.Split(key, " : ")
	return Queue{
		Group:   data[0],
		Subject: data[1],
	}
}

func (q *Queue) Key() string {
	return fmt.Sprintf("%s : %s", q.Group, q.Subject)
}

type QueueEntry struct {
	ChatID string
}
