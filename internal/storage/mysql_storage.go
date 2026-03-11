package storage

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/namezzy/gtodo/internal/model"
)

var _ Storage = (*MySQLStorage)(nil)

type MySQLStorage struct {
	db *sql.DB
}

func NewMySQLStorage(dsn string) (*MySQLStorage, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接 MySQL 失败: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("MySQL 连接测试失败: %w", err)
	}

	s := &MySQLStorage{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}
	return s, nil
}

func (s *MySQLStorage) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id          INT AUTO_INCREMENT PRIMARY KEY,
		description VARCHAR(500) NOT NULL,
		priority    VARCHAR(20)  NOT NULL DEFAULT 'medium',
		created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		done_at     DATETIME     NULL,
		status      VARCHAR(20)  NOT NULL DEFAULT 'todo'
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err := s.db.Exec(query)
	return err
}

func (s *MySQLStorage) Load() ([]model.Task, error) {
	rows, err := s.db.Query(
		"SELECT id, description, priority, created_at, done_at, status FROM tasks ORDER BY id ASC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		var doneAt sql.NullTime
		if err := rows.Scan(&t.ID, &t.Description, &t.Priority, &t.CreatedAt, &doneAt, &t.Status); err != nil {
			return nil, err
		}
		if doneAt.Valid {
			t.DoneAt = doneAt.Time
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (s *MySQLStorage) Save(tasks []model.Task) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM tasks"); err != nil {
		return err
	}

	for _, t := range tasks {
		var doneAt *time.Time
		if !t.DoneAt.IsZero() {
			doneAt = &t.DoneAt
		}
		_, err := tx.Exec(
			"INSERT INTO tasks (id, description, priority, created_at, done_at, status) VALUES (?, ?, ?, ?, ?, ?)",
			t.ID, t.Description, t.Priority, t.CreatedAt, doneAt, t.Status,
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *MySQLStorage) NextID(tasks []model.Task) int {
	if len(tasks) == 0 {
		return 1
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })
	return tasks[len(tasks)-1].ID + 1
}

func (s *MySQLStorage) AddTask(task model.Task) error {
	var doneAt *time.Time
	if !task.DoneAt.IsZero() {
		doneAt = &task.DoneAt
	}
	result, err := s.db.Exec(
		"INSERT INTO tasks (description, priority, created_at, done_at, status) VALUES (?, ?, ?, ?, ?)",
		task.Description, task.Priority, time.Now(), doneAt, task.Status,
	)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	task.ID = int(id)
	return nil
}

func (s *MySQLStorage) UpdateTask(task model.Task) error {
	var doneAt *time.Time
	if !task.DoneAt.IsZero() {
		doneAt = &task.DoneAt
	}
	result, err := s.db.Exec(
		"UPDATE tasks SET description=?, priority=?, done_at=?, status=? WHERE id=?",
		task.Description, task.Priority, doneAt, task.Status, task.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("找不到 ID 为 %d 的事项", task.ID)
	}
	return nil
}

func (s *MySQLStorage) DeleteTask(id int) error {
	result, err := s.db.Exec("DELETE FROM tasks WHERE id=?", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("找不到 ID 为 %d 的事项", id)
	}
	return nil
}

func (s *MySQLStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
