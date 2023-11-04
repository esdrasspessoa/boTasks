package repository

import (
	"botasks/internal/task"
	"errors"
	"sync"
)

type TaskRepository interface {
	Create(task task.Task) (string, error)
	GetByID(taskID string) (*task.Task, error)
	Update(task task.Task) error
	Delete(taskID string) error
}

type MemoryTask struct {
	ID          string
	Title       string
	Description string
}

type MemoryTaskRepository struct {
	tasks map[string]MemoryTask
	mu    sync.Mutex
}


func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{
		tasks: make(map[string]MemoryTask),
	}
}

func (r *MemoryTaskRepository) Create(task MemoryTask) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	taskID := task.ID
	r.tasks[taskID] = task
	return taskID, nil
}

func (r *MemoryTaskRepository) GetByID(taskID string) (*MemoryTask, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[taskID]
	if !exists {
		return nil, errors.New("Task not found")
	}
	return &task, nil
}

func (r *MemoryTaskRepository) GetTasksByIDs(taskIDs []string) ([]MemoryTask, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks := make([]MemoryTask, 0, len(taskIDs))
	for _, taskID := range taskIDs {
		task, exists := r.tasks[taskID]
		if exists {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func (r *MemoryTaskRepository) Update(task MemoryTask) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.tasks[task.ID]
	if !exists {
		return errors.New("Task not found")
	}

	r.tasks[task.ID] = task
	return nil
}

func (r *MemoryTaskRepository) Delete(taskID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.tasks[taskID]
	if !exists {
		return errors.New("Task not found")
	}

	delete(r.tasks, taskID)
	return nil
}
