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
	const op = "mirea.GetSubjects"

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "GET", group.ICalLink, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем ответ
	resp, err := c.cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	subjects, err := ical.NewDecoder(resp.Body).Decode()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return subjects, nil
}
