package bot

import (
	"errors"
	"fmt"

	"github.com/iskanye/mirea-queue/internal/models"
	"github.com/iskanye/mirea-queue/internal/services"
	tele "gopkg.in/telebot.v4"
)

func (b *Bot) Start(c tele.Context) error {
	user, err := b.usersService.GetUser(b.ctx, c.Chat().ID)
	if err == nil {
		// Пользователь найден - приветствуем его
		return c.Send(fmt.Sprintf("Привет %s из группы %s", user.Name, user.Group))
	}
	if !errors.Is(err, services.ErrNotFound) {
		return err
	}

	// Если пользователь не существует начинаем диалоговую цепочку
	return b.Dialogue(c, func(ch <-chan *tele.Message, c tele.Context) error {
		msg, err := c.Bot().Send(c.Chat(), "Введите группу")
		if err != nil {
			return err
		}

		groupMsg := <-ch
		err = c.Bot().Delete(groupMsg)
		if err != nil {
			return err
		}

		msg, err = c.Bot().Edit(msg, "Введите своё имя и фамилию")
		if err != nil {
			return err
		}

		usernameMsg := <-ch
		err = c.Bot().Delete(usernameMsg)
		if err != nil {
			return err
		}

		user := models.User{
			Name:  usernameMsg.Text,
			Group: groupMsg.Text,
		}

		user, err = b.usersService.CreateUser(b.ctx, c.Chat().ID, user)
		if err != nil {
			return err
		}

		msg, err = c.Bot().Edit(msg, "Успешно зарегистрированы")
		return err
	})
}

func (b *Bot) Edit(c tele.Context) error {
	_, err := b.usersService.GetUser(b.ctx, c.Chat().ID)
	if errors.Is(err, services.ErrNotFound) {
		return c.Send("Вы не зарегистрированы")
	}
	if err != nil {
		return err
	}

	return b.Dialogue(c, func(ch <-chan *tele.Message, c tele.Context) error {
		msg, err := c.Bot().Send(c.Chat(), "Введите группу")
		if err != nil {
			return err
		}

		groupMsg := <-ch
		err = c.Bot().Delete(groupMsg)
		if err != nil {
			return err
		}

		msg, err = c.Bot().Edit(msg, "Введите своё имя и фамилию")
		if err != nil {
			return err
		}

		usernameMsg := <-ch
		err = c.Bot().Delete(usernameMsg)
		if err != nil {
			return err
		}

		user := models.User{
			Name:  usernameMsg.Text,
			Group: groupMsg.Text,
		}

		user, err = b.usersService.UpdateUser(b.ctx, c.Chat().ID, user)
		if err != nil {
			return err
		}

		msg, err = c.Bot().Edit(msg, "Успешно изменены данные")
		return err
	})
}
