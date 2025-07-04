package db

import (
	"database/sql"
	"errors"
	"fmt"
)

type Task struct {
	ID      int    `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func AddTask(task *Task) (int64, error) {

	var id int64
	// определите запрос
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ? ,?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err = res.LastInsertId()
	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	var tasks []*Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
	res, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}

	defer res.Close()

	for res.Next() {
		var (
			id      int
			date    string
			title   string
			comment string
			repeat  string
		)

		err = res.Scan(&id, &date, &title, &comment, &repeat)
		if err != nil {
			return nil, err
		}
		task := &Task{
			ID:      id,
			Date:    date,
			Title:   title,
			Comment: comment,
			Repeat:  repeat,
		}

		tasks = append(tasks, task)
	}

	return tasks, err
}

func GetTask(id string) (*Task, error) {

	if id == "" {
		return nil, errors.New("некорректный ID")
	}

	task := Task{}
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id=?`
	res := db.QueryRow(query, id)

	err := res.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("задача не найдена")
		}
		return nil, err
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	// параметры пропущены, не забудьте указать WHERE

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("некорректный id")
	}
	return nil
}

func UpdateDate(next string, id string) error {

	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := db.Exec(query, next, id)
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("некорректный id")
	}
	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача с ID %s не найдена", id)
	}

	return nil
}
