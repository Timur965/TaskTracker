package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// Задача.
type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// Tasks возвращает список задач из БД.
func (s *Storage) tasks(taskID, authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2)
		ORDER BY id;
	`,
		taskID,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)

	}
	return tasks, rows.Err()
}

func (s *Storage) TasksOfAuthor(authorID int) ([]Task, error) {
	return s.tasks(0, authorID)
}

func (s *Storage) AllTasks() ([]Task, error) {
	return s.tasks(0, 0)
}

func (s *Storage) UpdateTask(taskID int, newTask Task) error {
	_, err := s.db.Exec(context.Background(), `
	UPDATE tasks
	SET 
		closed = $2
	WHERE
   		id = $1;
	`,
		taskID,
		newTask.Closed,
	)

	return err
}

func (s *Storage) DeleteTask(taskID int) error {
	_, err := s.db.Exec(context.Background(), `
	DELETE FROM tasks WHERE id = $1;
	`,
		taskID,
	)

	return err
}

func (s *Storage) TasksOfLabel(label string) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT tasks.id, opened, closed, author_id, assigned_id, title, content FROM tasks, labels, tasks_labels
		WHERE labels.name = $1 AND tasks.id = tasks_labels.task_id AND
		tasks_labels.label_id = labels.id
	`,
		label,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)

	}
	return tasks, rows.Err()
}

func (s *Storage) NewTask(t Task) error {
	_, err := s.db.Exec(context.Background(), `
		INSERT INTO tasks (author_id, assigned_id, title, content)
		VALUES ($1, $2, $3, $4);
		`,
		t.AuthorID,
		t.AssignedID,
		t.Title,
		t.Content,
	)
	return err
}
