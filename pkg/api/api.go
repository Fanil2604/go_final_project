package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const DateFormat = "20060102"

const secretKey = "secret123"

type Password struct {
	Password string `json:"password"`
}

type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.StandardClaims
}

// Регистрация обработчика для маршрута /api/nextdate
func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", Auth(taskHandler))
	http.HandleFunc("/api/tasks", Auth(tasksListHandler))
	http.HandleFunc("/api/task/done", Auth(tasksDoneHandler))
	http.HandleFunc("/api/signin", signinHandler)
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	// Логика обработчика
	//fmt.Fprintln(w, "Приветствуем вас на нашем сервере!")
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	// Парсинг текущей даты
	currentDate, err := time.Parse(DateFormat, now)
	if err != nil {
		fmt.Fprintln(w, "Ошибка при парсинге даты:", err)
		return
	}

	nextDate, err := NextDate(currentDate, date, repeat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "string")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(nextDate))
	if err != nil {
		log.Println("Error write in func GetNextDate:", err)
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// обработка других методов будет добавлена на следующих шагах
	case http.MethodPost:
		addTaskHandle(w, r)
	case http.MethodPut:
		updateTaskHandle(w, r)
	case http.MethodGet:
		getTaskHandle(w, r)
	case http.MethodDelete:
		deleteTaskHandle(w, r)
	}
}

func tasksListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// обработка других методов будет добавлена на следующих шагах
	case http.MethodGet:
		tasksHandler(w, r)
	}
}

func tasksDoneHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// обработка других методов будет добавлена на следующих шагах
	case http.MethodPost:
		taskDoneHandler(w, r)
	}
}
