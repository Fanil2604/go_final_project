package api

import (
	"encoding/json"
	"fmt"
	"go_final_project/pkg/db"
	"log"
	"net/http"
	"strconv"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		fmt.Println(err)
		return
	}

	if tasks == nil {
		tasks = []*db.Task{}
	}

	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

func getTaskHandle(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		//fmt.Println("не указан идентификатор")
		//writeERROR(w, "not id", http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	task, err := db.GetTask(id) // в параметре максимальное количество записей
	if err != nil {
		writeJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	writeJson(w, task)
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	_, err := strconv.Atoi(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "id is not a number"})
		return
	}

	t, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "can't get task"})
		return
	}

	if t.Repeat == "" {
		log.Println("Repeat is empty, task will delete")
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": "can't delete task"})
			return
		}
	}

	if t.Repeat != "" {
		t.Date, err = NextDate(time.Now(), t.Date, t.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": string(err.Error())})
			return
		}

		err = db.UpdateTask(t)
		if err != nil {
			writeJson(w, map[string]string{"error": string(err.Error())})
			return
		}
	}

	writeJson(w, map[string]string{})
}

func updateTaskHandle(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, jsonError("Не указан заголовок задачи"), http.StatusBadRequest)
		return
	}

	// Проверка формата даты
	if newTask.Date == "" {
		newTask.Date = time.Now().Format("20060102")
	} else {
		_, err = time.Parse("20060102", newTask.Date)
		if err != nil {
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
	err = db.UpdateTask(&newTask)
	if err != nil {
		http.Error(w, jsonError("Ошибка при добавлении задачи"), http.StatusInternalServerError)
		return
	}

	// Возврат идентификатора созданной задачи
	response := map[string]interface{}{}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

func writeJson(w http.ResponseWriter, data interface{}) {
	// Сериализуем данные в JSON, возвращаем ошибку, если она возникает:
	js, _ := json.Marshal(data)

	// Добавляем в JSON символ новой строки для удобства просмотра в терминале:
	js = append(js, '\n')
	// Устанавливаем заголовок Content-Type: application/json:
	w.Header().Set("Content-Type", "application/json")
	// Записываем JSON в ResponseWriter:
	w.Write(js)

}

func deleteTaskHandle(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		//fmt.Println("не указан идентификатор")
		//writeERROR(w, "not id", http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	err := db.DeleteTask(id) // в параметре максимальное количество записей
	if err != nil {
		writeJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	writeJson(w, map[string]string{})
}
