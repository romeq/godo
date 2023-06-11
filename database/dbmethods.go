package database

import (
	"github.com/romeq/godo/todos"
)

func GetTodos() ([]todos.Todo, error) {
	ensureDb()
	rows, err := db.Query(`SELECT * FROM todos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks todos.Todos
	for rows.Next() {
		var todo todos.Todo
		if err := rows.Scan(
			&todo.ID, &todo.Task, &todo.Done, &todo.Date, &todo.Deadline); err != nil {
			return tasks.SortByDone(), err
		}
		tasks = append(tasks, todo)
	}

	return tasks.SortByDone(), nil
}

func GetTodosAmount() (int32, error) {
	ensureDb()
	row := db.QueryRow("SELECT COUNT(*) FROM todos")
	if err := row.Err(); err != nil {
		return 0, err
	}

	var amount int32
	err := row.Scan(&amount)
	return amount, err
}

func NewTodo(task string, deadline int, date int) (int64, error) {
	ensureDb()

	x, err := db.Exec(`INSERT INTO todos(task, done, deadline, date) VALUES (?, ?, ?, ?)`,
		task, false, deadline, date)
	if err != nil {
		return 0, err
	}

	insertId, err := x.LastInsertId()
	return insertId, err
}

func UpdateDoneById(id int, done bool) error {
	ensureDb()
	_, err := db.Exec(`UPDATE todos SET done = ? WHERE id = ?`, done, id)
	return err
}

func RemoveTodoById(id int) error {
	ensureDb()
	_, err := db.Exec(`DELETE FROM todos WHERE id = ?`, id)
	return err
}
