package bot

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/iskanye/mirea-queue/internal/models"
	"github.com/iskanye/mirea-queue/internal/services"
	"gopkg.in/telebot.v4"
)

// Пушает в очередь
func (b *Bot) Push(c telebot.Context) error {
	err := b.Dialogue(c, func(ch <-chan *telebot.Message, c telebot.Context) error {
		msg, err := c.Bot().Send(c.Chat(), "Введите название учебной дисциплины")
		if err != nil {
			return err
		}

		subjectMsg := <-ch

		user := c.Get("user").(models.User)

		queue := models.Queue{
			Group:   user.Group,
			Subject: strings.TrimSpace(subjectMsg.Text),
		}

		entry := models.QueueEntry{
			ChatID: fmt.Sprint(c.Chat().ID),
		}

		pos, err := b.queueService.Push(b.ctx, queue, entry)
		if err != nil {
			if errors.Is(err, services.ErrAlreadyInQueue) {
				_, err := c.Bot().Edit(msg, "Вы уже в очереди")
				return err
			}
			return err
		}

		_, err = c.Bot().Edit(msg,
			fmt.Sprintf(
				"Очередь %s\nТекущая ваша позиция: %d",
				queue.Key(), pos,
			),
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Попает из очереди
func (b *Bot) Pop(c telebot.Context) error {
	err := b.Dialogue(c, func(ch <-chan *telebot.Message, c telebot.Context) error {
		msg, err := c.Bot().Send(c.Chat(), "Введите название учебной дисциплины")
		if err != nil {
			return err
		}

		subjectMsg := <-ch

		user := c.Get("user").(models.User)

		queue := models.Queue{
			Group:   user.Group,
			Subject: subjectMsg.Text,
		}

		entry, err := b.queueService.Pop(b.ctx, queue)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				_, err := c.Bot().Edit(msg, "Очередь не найдена")
				return err
			}
			return err
		}

		// Гарантировано что айди конвертируетсся в инт64
		chatID, _ := strconv.Atoi(entry.ChatID)

		user, err = b.usersService.GetUser(b.ctx, int64(chatID))
		if err != nil {
			return err
		}

		_, err = c.Bot().Edit(msg,
			fmt.Sprintf(
				"Очередь по предмету %s\nНа сдачу приглашается %s",
				queue.Subject, user.Name,
			),
		)
		if err != nil {
			return err
		}

		// Получаем чат того чья очередь
		chat, err := c.Bot().ChatByID(int64(chatID))
		if err != nil {
			return err
		}

		_, err = c.Bot().Send(chat,
			fmt.Sprintf(
				"Вы приглашаетесь на сдачу по предмету %s",
				queue.Subject,
			),
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Пропускает следующего в очереди
func (b *Bot) LetAhead(c telebot.Context) error {
	err := b.Dialogue(c, func(ch <-chan *telebot.Message, c telebot.Context) error {
		msg, err := c.Bot().Send(c.Chat(), "Введите название учебной дисциплины")
		if err != nil {
			return err
		}

		subjectMsg := <-ch

		user := c.Get("user").(models.User)

		queue := models.Queue{
			Group:   user.Group,
			Subject: subjectMsg.Text,
		}

		entry := models.QueueEntry{
			ChatID: fmt.Sprint(c.Chat().ID),
		}

		err = b.queueService.LetAhead(b.ctx, queue, entry)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				_, err := c.Bot().Edit(msg, "Очередь не найдена, либо вы в неё не записаны")
				return err
			}
			if errors.Is(err, services.ErrQueueEnd) {
				_, err := c.Bot().Edit(msg, "Вы последний в очереди")
				return err
			}
			return err
		}

		pos, err := b.queueService.Pos(b.ctx, queue, entry)
		if err != nil {
			return err
		}

		_, err = c.Bot().Edit(
			msg,
			fmt.Sprintf("Вы успешно пропустили следующего в очереди\nВаша текущая позиция в очереди - %d", pos),
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) ChooseSubject(c telebot.Context) error {
	err := b.Dialogue(c, func(ch <-chan *telebot.Message, c telebot.Context) error {
		msg, err := c.Bot().Send(c.Chat(), "Введите название учебной дисциплины")
		if err != nil {
			return err
		}

		subjectMsg := <-ch

		err = c.Bot().Delete(msg)
		if err != nil {
			return err
		}

		user := c.Get("user").(models.User)

		queue := models.Queue{
			Group:   user.Group,
			Subject: subjectMsg.Text,
		}

		entry := models.QueueEntry{
			ChatID: fmt.Sprint(c.Chat().ID),
		}

		pos, err := b.queueService.Pos(b.ctx, queue, entry)

		err = c.Edit(
			fmt.Sprintf("Ваша текущая позиция в очереди %s - %d", queue.Key(), pos),
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
