package bot

import (
	"errors"
	"fmt"

	"github.com/iskanye/mirea-queue/internal/models"
	"github.com/iskanye/mirea-queue/internal/services/users"
	"gopkg.in/telebot.v4"
)

func (b *Bot) Start(c telebot.Context) error {
	return b.Dialogue(c, func(ch <-chan string, c telebot.Context) error {
		user, err := b.usersService.GetUser(b.ctx, c.Chat().ID)
		if err == nil {
			return c.Send(fmt.Sprintf("Привет %s из группы %s", user.Name, user.Group))
		}
		if !errors.Is(err, users.ErrNotFound) {
			return err
		}

		err = c.Send("Введите группу")
		if err != nil {
			return err
		}

		group := <-ch

		err = c.Send("Введите своё имя и фамилию")
		if err != nil {
			return err
		}

		username := <-ch

		user = models.User{
			Name:  username,
			Group: group,
		}

		user, err = b.usersService.CreateUser(b.ctx, c.Chat().ID, user)
		if err != nil {
			return err
		}

		return c.Send("Успешно зарегистрированы")
	})
}
