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

func (b *Bot) ChooseGroup(c tele.Context) error {
	ch := b.channels[c.Chat().ID]
	ch <- c.Data()
	close(ch)
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

// Вернуться на главную страницу
func (b *Bot) Return(c tele.Context) error {
	user := c.Get("user").(models.User)

	err := c.Delete()
	if err != nil {
		return err
	}

	return b.showProfile(c, user)
}

// Функция получения пользователя из ввода
func (b *Bot) getUser(c tele.Context) (models.User, error) {
	var user models.User

	msg, err := c.Bot().Send(c.Chat(), "Введите группу")
	if err != nil {
		return models.User{}, err
	}

	ch := make(chan string, 1)
	b.channels[c.Chat().ID] = ch

	// Позволяем указывать только актуальную студенческую группу
	var group string
	for i := range ch {
		groups, err := b.scheduleService.GetGroups(b.ctx, i)
		// Группа не найдена в расписании
		if errors.Is(err, services.ErrNotFound) {
			msg, err = c.Bot().Edit(msg, "Данная группа не найдена. Попробуйте ещё раз")
			continue
		} else if err != nil {
			return models.User{}, err
		}

		// Если сразу получили группу то можем не продолжать
		if len(groups) == 1 {
			group = groups[0].Name
			break
		}

		// Создаю кнопки под сообщением
		groupMarkup := &tele.ReplyMarkup{}
		btns := make([]tele.Btn, len(groups))
		for j := range groups {
			btns[j] = groupMarkup.Data(groups[j].Name, b.groupBtnUnique, groups[j].Name)
		}
		groupMarkup.Inline(
			groupMarkup.Split(1, btns)...,
		)

		msg, err = c.Bot().Edit(msg, "Выберите группу", groupMarkup)
		if err != nil && !errors.Is(err, tele.ErrSameMessageContent) {
			return models.User{}, err
		}
	}

	// Читаем оставшиеся данные
	err = b.Dialogue(c, func(ch <-chan string, c tele.Context) error {
		msg, err = c.Bot().Edit(msg, "Введите своё имя и фамилию")
		if err != nil {
			return err
		}

		username := <-ch

		msg, err = c.Bot().Edit(msg, "Введите токен админа(если есть)")
		if err != nil {
			return err
		}

		token := <-ch

		err = c.Bot().Delete(msg)
		if err != nil {
			return err
		}

		user = models.User{
			Name:        strings.TrimSpace(username),
			Group:       strings.TrimSpace(group),
			QueueAccess: b.adminService.ValidateToken(strings.TrimSpace(token)),
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
