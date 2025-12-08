package schedule

import (
	"net/http"

	"github.com/iskanye/mirea-queue/internal/interfaces"
)

const scheduleUrl = "https://schedule-of.mirea.ru/schedule/api/"

// Проверка на реализацию интерфейсов, чтобы ничего не поломалось
var (
	_ interfaces.GroupProvider    = (*Client)(nil)
	_ interfaces.ScheduleProvider = (*Client)(nil)
)

type Client struct {
	cl *http.Client
}

func New() *Client {
	return &Client{
		cl: http.DefaultClient,
	}
}
