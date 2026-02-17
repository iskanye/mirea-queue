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

// Обновляет данные очереди
func (b *Bot) Refresh(c telebot.Context) error {
	queue := c.Get("queue").(models.Queue)

	entry := models.QueueEntry{
		ChatID: fmt.Sprint(c.Chat().ID),
	}

	return b.showSubject(c, queue, entry)
}

// Пушает в очередь
func (b *Bot) Push(c telebot.Context) error {
	queue := c.Get("queue").(models.Queue)

	entry := models.QueueEntry{
		ChatID: fmt.Sprint(c.Chat().ID),
	}

	if err := b.queueService.Push(b.ctx, queue, entry); err != nil {
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
	queue := c.Get("queue").(models.Queue)

	entry := models.QueueEntry{
		ChatID: fmt.Sprint(c.Chat().ID),
	}

	err := b.queueService.LetAhead(b.ctx, queue, entry)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			return c.Send("Вы не записаны в очередь")
		}
		if errors.Is(err, services.ErrQueueEnd) {
			return c.Send("Вы последний в очереди")
		}
		return err
	}

	return b.showSubject(c, queue, entry)
}

// Выбрать предмет
func (b *Bot) ChooseSubject(c telebot.Context) error {
	user := c.Get("user").(models.User)
	groups, err := b.scheduleService.GetGroups(b.ctx, user.Group)
	if err != nil {
		return err
	}

	// Группа гарантировано будет одна
	subjects, err := b.scheduleService.GetSubjects(b.ctx, groups[0])
	if err != nil {
		return err
	}

	// Создаю кнопки под сообщением
	subjectMarkup := &telebot.ReplyMarkup{}
	btns := make([]telebot.Btn, len(subjects))
	var btnText strings.Builder
	for i := range subjects {
		// В качестве полезной нагрузки возьмём первое слово названия дисциплины
		// TODO: #17 Придумать способ хранения callback_data по-лучше для кнопок
		data, _, _ := strings.Cut(subjects[i], " ")

		queue := models.Queue{
			Group:   user.Group,
			Subject: data,
		}

		// Проверяем, есть ли уже очередь по этому предмету
		_, err := b.queueService.Range(b.ctx, queue)
		if err == nil {
			btnText.WriteString("🟩 ")
		} else if errors.Is(err, services.ErrNotFound) {
			btnText.WriteString("🟥 ")
		}
		btnText.WriteString(subjects[i])

		btns[i] = subjectMarkup.Data(btnText.String(), b.subjectBtnUnique, data)
		btnText.Reset()
	}
	subjectMarkup.Inline(
		subjectMarkup.Split(1, btns)...,
	)

	err = c.Edit("Выберите учебную дисциплину", subjectMarkup)
	if err != nil {
		return err
	}

	// Получаем название дисциплины
	ch := make(chan string, 1)
	b.channels[c.Chat().ID] = ch
	subject := <-ch
	close(ch)
	delete(b.channels, c.Chat().ID)

	queue := models.Queue{
		Group:   user.Group,
		Subject: subject,
	}

	entry := models.QueueEntry{
		ChatID: fmt.Sprint(c.Chat().ID),
	}

	err = b.queueService.SaveToCache(b.ctx, c.Chat().ID, queue)
	if err != nil {
		return err
	}

	return b.showSubject(c, queue, entry)
}

// Обработчик кнопки выбора предмета
func (b *Bot) ChooseSubjectButton(c telebot.Context) error {
	b.channels[c.Chat().ID] <- c.Data()
	return nil
}

// Очищает очередь
func (b *Bot) Clear(c telebot.Context) error {
	queue := c.Get("queue").(models.Queue)
	entry := models.QueueEntry{
		ChatID: fmt.Sprint(c.Chat().ID),
	}

	err := b.queueService.Clear(b.ctx, queue)
	if err != nil {
		return err
	}

	return b.showSubject(c, queue, entry)
}

// Удаляет пользователя из очереди
func (b *Bot) Remove(c telebot.Context) error {
	queue := c.Get("queue").(models.Queue)
	entry := models.QueueEntry{
		ChatID: fmt.Sprint(c.Chat().ID),
	}

	err := b.queueService.Remove(b.ctx, queue, entry)
	if err != nil {
		return err
	}

	return b.showSubject(c, queue, entry)
}

// Выводит на экран информацию об очереди
func (b *Bot) showSubject(
	c telebot.Context,
	queue models.Queue,
	entry models.QueueEntry,
) error {
	var sb strings.Builder
	sb.WriteString("Очередь " + queue.Key())

	entries, err := b.queueService.Range(b.ctx, queue)
	if errors.Is(err, services.ErrNotFound) {
		sb.WriteString("\nОчередь пуста")
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

			fmt.Fprintf(&sb, "\n%3d.  %s", i+1, user.Name)
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

	menu := b.subjectMenu
	if user := c.Get("user").(models.User); user.QueueAccess {
		menu = b.subjectAdminMenu
	}

	err = c.Edit(sb.String(), menu)
	if err != nil && !errors.Is(err, telebot.ErrSameMessageContent) {
		return err
	}

	return nil
}
