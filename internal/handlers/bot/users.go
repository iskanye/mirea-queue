package bot

import (
	"errors"
	"fmt"
	"strings"

	"github.com/iskanye/mirea-queue/internal/models"
	"github.com/iskanye/mirea-queue/internal/services"
	tele "gopkg.in/telebot.v4"
)

func (b *Bot) Start(c tele.Context) error {
	user, err := b.usersService.GetUser(b.ctx, c.Chat().ID)
	if err == nil {
		// Пользователь найден - приветствуем его
		return b.showProfile(c, user)
	}
	if !errors.Is(err, services.ErrNotFound) {
		return err
	}

	// Если пользователь не существует получаем его данные
	user, err = b.getUser(c)
	if err != nil {
		return err
	}

	user, err = b.usersService.CreateUser(b.ctx, c.Chat().ID, user)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) Edit(c tele.Context) error {
	// Получаем новые данные пользователя
	user, err := b.getUser(c)
	if err != nil {
		return err
	}

	user, err = b.usersService.UpdateUser(b.ctx, c.Chat().ID, user)
	if err != nil {
		return err
	}

	return c.Delete()
}

// Функция получения пользователя из ввода
func (b *Bot) getUser(c tele.Context) (models.User, error) {
	var user models.User
	err := b.Dialogue(c, func(ch <-chan *tele.Message, c tele.Context) error {
		msg, err := c.Bot().Send(c.Chat(), "Введите группу")
		if err != nil {
			return err
		}

		groupMsg := <-ch

		msg, err = c.Bot().Edit(msg, "Введите своё имя и фамилию")
		if err != nil {
			return err
		}

		usernameMsg := <-ch

		msg, err = c.Bot().Edit(msg, "Введите токен админа(если есть)")
		if err != nil {
			return err
		}

		tokenMsg := <-ch

		err = c.Bot().Delete(msg)
		if err != nil {
			return err
		}

		user = models.User{
			Name:        strings.TrimSpace(usernameMsg.Text),
			Group:       strings.TrimSpace(groupMsg.Text),
			QueueAccess: b.adminService.ValidateToken(strings.TrimSpace(tokenMsg.Text)),
		}

		return b.showProfile(c, user)
	})
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// Функция отображения профиля
func (b *Bot) showProfile(c tele.Context, user models.User) error {
	return c.Send(fmt.Sprintf(
		"Группа: %s\nФИО: %s\nПрава админа: %t",
		user.Group, user.Name, user.QueueAccess,
	), b.startMenu)
}
