package schedule

import "net/http"

const scheduleUrl = "https://schedule-of.mirea.ru/schedule/api/"

type Client struct {
	cl *http.Client
}

func New() *Client {
	return &Client{
		cl: &http.Client{},
	}
}
