package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/namezzy/gtodo/internal/model"
)

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

func (s *JSONStorage) Load() ([]model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

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
	// 按 ID 排序
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })
	return tasks, nil
}

func (s *JSONStorage) Save(tasks []model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *JSONStorage) NextID(tasks []model.Task) int {
	if len(tasks) == 0 {
		return 1
	}
	return tasks[len(tasks)-1].ID + 1
}
