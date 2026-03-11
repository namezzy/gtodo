package storage

import "github.com/namezzy/gtodo/internal/model"

// Storage 定义了任务存储的统一接口
type Storage interface {
	Load() ([]model.Task, error)
	Save(tasks []model.Task) error
	NextID(tasks []model.Task) int
	AddTask(task model.Task) error
	UpdateTask(task model.Task) error
	DeleteTask(id int) error
	Close() error
}
