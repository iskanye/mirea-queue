package ical_test

import (
	"os"
	"testing"

	"github.com/iskanye/mirea-queue/internal/lib/ical"
	"github.com/stretchr/testify/assert"
)

const testFilePath = "./test.ics"

var testSubjects = []string{
	"Архитектура вычислительных машин и систем",
	"Вычислительная математика",
	"Иностранный язык",
	"Конфигурационное управление",
	"Математический анализ",
	"Настройка и администрирование сервисного программного обеспечения",
	"Программирование на языке Джава",
	"Системы искусственного интеллекта и большие данные",
	"Структуры и алгоритмы обработки данных",
	"Физическая культура и спорт (элективная дисциплина)",
	"Фронтенд-разработка",
	"Экономическая культура",
}

func TestICal(t *testing.T) {
	t.Parallel()

	reader, _ := os.Open(testFilePath)
	subjects, _ := ical.NewDecoder(reader).Decode()

	for _, subject := range subjects {
		assert.Contains(t, testSubjects, subject)
	}
}
