package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"log"
	"net/http"
	"time"
)

func addTaskHandle(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос: %s %s", r.Method, r.URL.Path)

	var newTask db.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		log.Printf("Ошибка десериализации JSON: %v", err)
		http.Error(w, jsonError("Ошибка десериализации JSON"), http.StatusBadRequest)
		return
	}

	// Проверка обязательного поля title
	if newTask.Title == "" {
		log.Printf("Не указан заголовок задачи: %v", err)
		http.Error(w, jsonError("Не указан заголовок задачи"), http.StatusBadRequest)
		return
	}

	// Проверка формата даты
	if newTask.Date == "" {
		newTask.Date = time.Now().Format("20060102")
	} else {
		_, err = time.Parse("20060102", newTask.Date)
		if err != nil {
			log.Printf("Дата представлена в неправильном формате: %v", err)
			http.Error(w, jsonError("Дата представлена в неправильном формате"), http.StatusBadRequest)
			return
		}
	}

	// Проверка даты
	today := time.Now().Format("20060102")
	if newTask.Date < today {
		if newTask.Repeat == "" {
			newTask.Date = today
		} else {
			currentDate, _ := time.Parse("20060102", today)
			for {

				nextDate, err := NextDate(currentDate, newTask.Date, newTask.Repeat)
				newTask.Date = nextDate
				if err != nil {
					log.Printf("Неправильный формат правила повторения: %v", err)
					http.Error(w, jsonError("Неправильный формат правила повторения"), http.StatusBadRequest)
					return
				}
				if afterNow(nextDate, today) {

					break
				}
			}

		}
	}

	// Сохранение задачи в БД
	id, err := db.AddTask(&newTask)
	if err != nil {
		log.Printf("Ошибка при добавлении задачи: %v", err)
		http.Error(w, jsonError("Ошибка при добавлении задачи"), http.StatusInternalServerError)
		return
	}

	// Возврат идентификатора созданной задачи
	response := map[string]interface{}{"id": id}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

// jsonError формирует JSON-ответ с ошибкой
func jsonError(message string) string {
	errorResponse, _ := json.Marshal(map[string]string{"error": message})
	return string(errorResponse)
}
func afterNow(date string, now string) bool {
	currentDate, _ := time.Parse("20060102", date)
	nowDate, _ := time.Parse("20060102", now)
	return currentDate.After(nowDate)
}
