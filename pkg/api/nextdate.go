package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Парсинг текущей даты
	currentDate, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}

	// Логика для разных правил повторения
	switch repeat {
	case "y":
		// Если правило y, то добавляем год
		nextDate := currentDate.AddDate(1, 0, 0)

		return nextDate.Format(DateFormat), nil
	default:
		if strings.HasPrefix(repeat, "d") {
			parts := strings.Split(repeat, " ")
			if len(parts) == 2 {
				days, err := strconv.Atoi(parts[1])
				if err == nil {
					nextDate := currentDate.Add(time.Duration(days) * 24 * time.Hour)
					return nextDate.Format(DateFormat), nil
				}
			}
		}
		// Если не указано правило или оно неверно
		return "", fmt.Errorf("неверное правило повторения: %s", repeat)
	}

}
