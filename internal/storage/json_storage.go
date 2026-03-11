package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/namezzy/gtodo/internal/model"
)

// 确保 JSONStorage 实现 Storage 接口
var _ Storage = (*JSONStorage)(nil)

type JSONStorage struct {
	path string
	mu   sync.Mutex
}

func NewJSONStorage() (*JSONStorage, error) {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".gtodo", "tasks.json")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	return &JSONStorage{path: path}, nil
}

func (s *JSONStorage) load() ([]model.Task, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Task{}, nil
		}
		return nil, err
	}
	var tasks []model.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })
	return tasks, nil
}

func (s *JSONStorage) save(tasks []model.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *JSONStorage) Load() ([]model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.load()
}

func (s *JSONStorage) Save(tasks []model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.save(tasks)
}

func (s *JSONStorage) NextID(tasks []model.Task) int {
	if len(tasks) == 0 {
		return 1
	}
	return tasks[len(tasks)-1].ID + 1
}

func (s *JSONStorage) AddTask(task model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, err := s.load()
	if err != nil {
		return err
	}

	task.ID = s.NextID(tasks)
	task.CreatedAt = time.Now()
	tasks = append(tasks, task)
	return s.save(tasks)
}

func (s *JSONStorage) UpdateTask(task model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, err := s.load()
	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].ID == task.ID {
			tasks[i] = task
			return s.save(tasks)
		}
	}
	return fmt.Errorf("找不到 ID 为 %d 的事项", task.ID)
}

func (s *JSONStorage) DeleteTask(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, err := s.load()
	if err != nil {
		return err
	}

	newTasks := make([]model.Task, 0, len(tasks))
	found := false
	for _, t := range tasks {
		if t.ID == id {
			found = true
			continue
		}
		newTasks = append(newTasks, t)
	}

	if !found {
		return fmt.Errorf("找不到 ID 为 %d 的事项", id)
	}
	return s.save(newTasks)
}

func (s *JSONStorage) Close() error {
	return nil
}
