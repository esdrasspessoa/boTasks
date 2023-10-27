package repository

import (
	"botasks/internal/list"
	"botasks/internal/task"
)

type TaskRepository interface {
	Create(task task.Task) (string, error)
	GetByID(taskID string) (*task.Task, error)
	Update(task task.Task) error
	Delete(taskID string) error
}

type TaskListRepository interface {
	Create(taskList list.TaskList) (string, error)
	GetByID(taskListID string) (*list.TaskList, error)
	Update(taskList list.TaskList) error
	Delete(taskListID string) error
	AddTaskToList(taskID, taskListID string) error
	GetTasksByList(taskListID string) ([]list.Task, error)
}
