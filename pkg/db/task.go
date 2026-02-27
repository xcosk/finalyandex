package db

import (
	"fmt"
	"strconv"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	res, err := DB.Exec(
		"INSERT INTO scheduler(date, title, comment, repeat) VALUES(?, ?, ?, ?)",
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(limit int) ([]*Task, error) {
	rows, err := DB.Query(
		"SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		t := &Task{}
		var id int64
		if err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		t.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = make([]*Task, 0)
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	row := DB.QueryRow(
		"SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?",
		id,
	)

	var (
		task Task
		rid  int64
	)
	if err := row.Scan(&rid, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return nil, err
	}
	task.ID = strconv.FormatInt(rid, 10)
	return &task, nil
}

func UpdateTask(task *Task) error {
	res, err := DB.Exec(
		"UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for deleting task")
	}
	return nil
}

func UpdateDate(next string, id string) error {
	res, err := DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", next, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating date")
	}
	return nil
}
