package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"` // и т.д.
}

func AddTask(task *Task) (int64, error) {
	var id int64
	// определите запрос
	res, err := DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return 0, err
	}
	id, err = res.LastInsertId()
	if err != nil {
		log.Printf("Ошибка получения ID: %s\n", err)
	}
	return id, nil
}

func Tasks(limit int) ([]*Task, error) {
	rows, err := DB.Query(
		`SELECT * 
		 FROM scheduler
		 ORDER BY date ASC
		 LIMIT :limit`,
		sql.Named("limit", limit),
	)
	if err != nil {
		log.Println("can't get tasks by GetList:", err)
		return []*Task{}, err
	}

	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		t := Task{}

		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			log.Println("can't get tasks by GetList:", err)
			return []*Task{}, err
		}

		tasks = append(tasks, &t)
	}

	if rows.Err() != nil {
		log.Println("can't get tasks by GetList:", rows.Err())
		return []*Task{}, rows.Err()
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	task := &Task{}
	err := DB.QueryRow(
		"SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?",
		id,
	).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		return task, err
	}

	return task, nil
}

func DeleteTask(id string) error {
	task := &Task{}
	err := DB.QueryRow(
		"SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?",
		id,
	).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		return err
	}

	_, err = DB.Exec(
		`DELETE
		FROM scheduler
		WHERE id = :id`,
		sql.Named("id", id),
	)

	return err
}

func UpdateTask(task *Task) error {
	// параметры пропущены, не забудьте указать WHERE
	res, err := DB.Exec(
		`UPDATE scheduler 
		SET date=:date, title= :title, comment= :comment, repeat= :repeat
		WHERE id= :id`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)
	if err != nil {
		log.Println("can't update task:", err)
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}
