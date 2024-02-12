package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"time"
	"to-do-app/internal/http-server/handlers/todo"
)

type Storage struct {
	db *sql.DB
}

type Task struct {
	Id        int           `json:"id"`
	Task      string        `json:"task,omitempty"`
	CreatedAt time.Duration `json:"created_at,omitempty"`
	Active    bool          `json:"active,omitempty"`
	Status    string        `json:"status,omitempty"`
}

func New(host, port, user, password, dbName string) (*Storage, error) {
	const op = "storage.postgres.New"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	storage := &Storage{db: db}

	err = goose.Up(storage.db, "db/migrations")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return storage, nil
}

func (s *Storage) Create(task string) (int, error) {
	const op = "storage.postgres.Create"

	var id int

	query := `
INSERT INTO todo_list (task) VALUES ($1) RETURNING id;
`

	err := s.db.QueryRow(query, task).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) Delete(id int) error {
	const op = "storage.postgres.Delete"

	query := `
DELETE FROM todo_list WHERE id = $1;
`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Complete(id int) error {
	const op = "storage.postgres.Complete"

	query := `
UPDATE todo_list SET active = FALSE WHERE id = $1;
`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Edit(id int, editedTask string) error {
	const op = "storage.postgres.Edit"

	query := `
UPDATE todo_list SET task = $2 WHERE id = $1;
`

	err := s.db.QueryRow(query, id, editedTask)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetAll() ([]todo.TaskList, error) {
	const op = "storage.postgres.GetAll"

	var tasks []todo.TaskList

	rows, err := s.db.Query(`SELECT id, task, active, created_at FROM todo_list`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var t todo.TaskList
		if err := rows.Scan(&t.Id, &t.Task, &t.Active, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}
