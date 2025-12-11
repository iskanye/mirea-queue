package ical

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	PractisePrefix    = "SUMMARY:ПР " // Префикс практики
	LecturePrefix     = "SUMMARY:ЛК " // Префикс лекции
	IndependentPrefix = "SUMMARY:СР " // Префикс самостоятельной работы

	PrefixLen = 11 // Длина префиксов
)

type Decoder struct {
	s *bufio.Scanner
}

// Создаёт декодер ical файла
func NewDecoder(r io.Reader) *Decoder {
	scanner := bufio.NewScanner(r)
	return &Decoder{
		s: scanner,
	}
}

// Записывает найденные предметы в срез.
func (d *Decoder) Decode(v []string) error {
	const op = "ical.Decode"

	subjects := make(map[string]struct{})
	isSubject := false
	subject := strings.Builder{}

	for d.s.Scan() {
		line := d.s.Text()
		if strings.HasPrefix(line, PractisePrefix) ||
			strings.HasPrefix(line, LecturePrefix) ||
			strings.HasPrefix(line, IndependentPrefix) {
			// Если строка имеет префикс предмета записываем ее
			isSubject = true
			subject.WriteString(line[PrefixLen-1:])
		} else if isSubject && strings.HasPrefix(line, " ") {
			// Сами предметы могут состоять из множества линий,
			// об этом свидетельствует пробел
			subject.WriteString(line[1:])
		} else if isSubject {
			// Строка предмета создана - записываем её в множество
			isSubject = false
			subjects[subject.String()] = struct{}{}
			subject.Reset()
		}
	}

	if err := d.s.Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Записываем данные из множества в срез
	v = make([]string, 0, len(subjects))
	for k := range subjects {
		v = append(v, k)
	}

	return nil
}
