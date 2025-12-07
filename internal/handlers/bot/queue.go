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
	queue := c.Get("queue").(models.Queue)

	entry := models.QueueEntry{
		ChatID: fmt.Sprint(c.Chat().ID),
	}

	_, err := b.queueService.Push(b.ctx, queue, entry)
	if err != nil {
		if errors.Is(err, services.ErrAlreadyInQueue) {
			return c.Send("Вы уже в очереди")
		}
		return err
	}

	return b.showSubject(c, queue, entry)
}

// Попает из очереди
func (b *Bot) Pop(c telebot.Context) error {
	queue := c.Get("queue").(models.Queue)

	// Попаем челика из очереди
	entry, err := b.queueService.Pop(b.ctx, queue)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			return c.Send("Очередь пуста")
		}
		return err
	}

	// Айдишник гарантировано имеет тип int64, зуб даю
	chatID, _ := strconv.ParseInt(entry.ChatID, 10, 64)

	// Получаем бедолагу, которого только что попнули
	user, err := b.usersService.GetUser(b.ctx, chatID)
	if err != nil {
		return err
	}

	if chatID != c.Chat().ID {
		err = c.Send(fmt.Sprintf("На сдачу приглашается %s", user.Name))
		if err != nil {
			return err
		}
	}

	// Получаем чат того, кто щас сдавать пойдёт
	chat, err := c.Bot().ChatByID(chatID)
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

	// Обновляем информацию о положении в очереди
	// текущего пользователя (а не того которого мы попнули)
	entry = models.QueueEntry{
		ChatID: fmt.Sprint(c.Chat().ID),
	}

	return b.showSubject(c, queue, entry)
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

// Выбрать предмет
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

		err = b.queueService.SaveToCache(b.ctx, c.Chat().ID, queue)
		if err != nil {
			return err
		}

		return b.showSubject(c, queue, entry)
	})
	if err != nil {
		return err
	}

	return nil
}

// Выводит на экран информацию об очереди
func (b *Bot) showSubject(
	c telebot.Context,
	queue models.Queue,
	entry models.QueueEntry,
) error {
	var sb strings.Builder
	sb.WriteString(queue.Key())

	entries, err := b.queueService.Range(b.ctx, queue)
	if errors.Is(err, services.ErrNotFound) {
		sb.WriteString("\nОчередь не найдена")
	} else if err == nil {
		// Находим имена пользователей
		for i, entry := range entries {
			chatID, err := strconv.ParseInt(entry.ChatID, 10, 64)
			if err != nil {
				return err
			}

			user, err := b.usersService.GetUser(b.ctx, chatID)
			if err != nil {
				return err
			}

			sb.WriteString(fmt.Sprintf("\n%d: %s", i+1, user.Name))
		}

		// Находим позицию текущего пользователя
		pos, err := b.queueService.Pos(b.ctx, queue, entry)

		msgText := fmt.Sprintf("\nВаша текущая позиция в очереди - %d", pos)
		if errors.Is(err, services.ErrNotFound) {
			msgText = "\nВы не записаны в очередь"
		} else if err != nil {
			return err
		}

		sb.WriteString(msgText)
	} else {
		return err
	}

	err = c.Edit(sb.String(), b.subjectMenu)
	if err != nil && !errors.Is(err, telebot.ErrSameMessageContent) {
		return err
	}

	return nil
}
