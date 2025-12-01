package bot

import (
	"fmt"

	"github.com/iskanye/mirea-queue/internal/models"
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
			Subject: subjectMsg.Text,
		}

		entry := models.QueueEntry{
			Student: user.Name,
		}

		pos, err := b.queueService.Push(b.ctx, queue, entry)
		if err != nil {
			return err
		}

		_, err = c.Bot().Edit(msg,
			fmt.Sprintf(
				"Очередь %s:%s\nТекущая ваша позиция: %d",
				queue.Group, queue.Subject, pos,
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
			return err
		}

		_, err = c.Bot().Edit(msg,
			fmt.Sprintf(
				"Очередь %s:%s\nНа сдачу приглашается %s",
				queue.Group, queue.Subject, entry.Student,
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
