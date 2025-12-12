package schedule

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/iskanye/mirea-queue/internal/models"
)

type groupResponse struct {
	Groups []models.Group `json:"data"`
}

func (c *Client) GetGroups(
	ctx context.Context,
	group string,
	limit int,
) ([]models.Group, error) {
	const op = "mirea.GetGroups"

	// Добавляем параметры в запрос
	query := url.Values{}
	query.Set("match", group)
	query.Set("limit", fmt.Sprint(limit))

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "GET", scheduleUrl+"search?"+query.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем ответ
	resp, err := c.cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	// Парсим ответ (по хорошему должен быть жсон файлик)
	var data groupResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if len(data.Groups) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrNotFound)
	}

	return data.Groups, nil
}
