package bot

import (
	"github.com/iskanye/mirea-queue/internal/models"
	"gopkg.in/telebot.v4"
)

func (b *Bot) Start(c telebot.Context) error {
	bot := c.Bot()

	groupMsg, err := bot.Send(c.Sender(), "Введите свою группу")
	if err != nil {
		return err
	}

	group := groupMsg.Text

	usernameMsg, err := bot.Send(c.Sender(), "Введите своё имя и фамилию")
	if err != nil {
		return err
	}

	username := usernameMsg.Text

	user := models.User{
		Name:  username,
		Group: group,
	}

	user, err = b.usersService.CreateUser(b.ctx, c.Chat().ID, user)
	if err != nil {
		return err
	}

	return c.Send("Успешно")
}
