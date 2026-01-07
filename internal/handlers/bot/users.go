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

	return b.showProfile(c, user)
}

func (b *Bot) ChooseGroup(c tele.Context) error {
	ch := b.channels[c.Chat().ID]
	ch <- c.Data()
	close(ch)
	return nil
}

func (b *Bot) Edit(c tele.Context) error {
	err := c.Delete()
	if err != nil {
		return err
	}

	// Получаем новые данные пользователя
	user, err := b.getUser(c)
	if err != nil {
		return err
	}

	user, err = b.usersService.UpdateUser(b.ctx, c.Chat().ID, user)
	if err != nil {
		return err
	}

	return b.showProfile(c, user)
}

// Вернуться на главную страницу
func (b *Bot) Return(c tele.Context) error {
	user := c.Get("user").(models.User)
	c.Set("msg", c.Message())
	return b.showProfile(c, user)
}

// Функция получения пользователя из ввода
func (b *Bot) getUser(c tele.Context) (models.User, error) {
	// Получаем данные юзера, если они есть значит что пользователь
	// меняет свои данные, в переменной ok храним меняет ли пользователь данные
	user, ok := c.Get("user").(models.User)

	// Если данные изменяют добавляем опцию вернуть те же данные
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(
		menu.Row(menu.Text(user.Group)),
	)

	var err error
	if ok {
		err = c.Send("Введите группу", menu)
	} else {
		err = c.Send("Введите группу")
	}
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
			err = c.Send("Данная группа не найдена. Попробуйте ещё раз")
			if err != nil {
				return models.User{}, err
			}
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

		err = c.Send("Выберите группу", groupMarkup)
		if err != nil && !errors.Is(err, tele.ErrSameMessageContent) {
			return models.User{}, err
		}
	}

	// Читаем оставшиеся данные
	err = b.dialogue(c, func(ch <-chan string, c tele.Context) error {
		var err error

		menu.Reply(
			menu.Row(menu.Text(user.Name)),
		)

		if ok {
			err = c.Send("Введите своё имя и фамилию", menu)
		} else {
			err = c.Send("Введите своё имя и фамилию")
		}
		if err != nil {
			return err
		}

		username := <-ch

		noMenu := &tele.ReplyMarkup{RemoveKeyboard: true}

		err = c.Send("Введите токен админа(если есть)", noMenu)
		if err != nil {
			return err
		}

		token := <-ch

		user = models.User{
			Name:        strings.TrimSpace(username),
			Group:       strings.TrimSpace(group),
			QueueAccess: b.adminService.ValidateToken(strings.TrimSpace(token)),
		}

		return nil
	})
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// Функция отображения профиля
func (b *Bot) showProfile(c tele.Context, user models.User) error {
	profile := fmt.Sprintf(
		"Группа: %s\nФИО: %s\nПрава админа: %t",
		user.Group, user.Name, user.QueueAccess,
	)
	if msg, ok := c.Get("msg").(tele.Editable); ok {
		_, err := c.Bot().Edit(msg, profile, b.startMenu)
		return err
	} else {
		return c.Send(profile, b.startMenu)
	}
}
