package schedule

import (
	"context"
	"fmt"
	"net/http"

	"github.com/iskanye/mirea-queue/internal/lib/ical"
	"github.com/iskanye/mirea-queue/internal/models"
)

func (c *Client) GetSubjects(
	ctx context.Context,
	group models.Group,
) ([]string, error) {
	const op = "schedule.GetSubjects"

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "GET", scheduleUrl+"ical/1/"+fmt.Sprint(group.ID), nil)
	if err == nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем ответ
	resp, err := c.cl.Do(req)
	if err == nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	var subjects []string
	err = ical.NewDecoder(resp.Body).Decode(subjects)
	if err == nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return subjects, nil
}
