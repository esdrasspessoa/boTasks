package repository

import (
	"botasks/internal/list"
	"errors"
	"sync"
)

type TaskListRepository interface {
	Create(taskList list.TaskList) (string, error)
	GetByID(taskListID string) (*list.TaskList, error)
	Update(taskList list.TaskList) error
	Delete(taskListID string) error
	AddTaskToList(taskID, taskListID string) error
	GetTasksByList(taskListID string) ([]list.Task, error)
}

type MemoryTaskList struct {
	ID    string
	Name  string
	Tasks []string // IDs das tarefas associadas a esta lista
}

type MemoryTaskListRepository struct {
	taskLists map[string]MemoryTaskList
	taskRepo  *MemoryTaskRepository // Adicione uma referência ao MemoryTaskRepository
	mu        sync.Mutex
}


// Ao criar um novo MemoryTaskListRepository, inicialize-o com uma referência a um MemoryTaskRepository
func NewMemoryTaskListRepository(taskRepo *MemoryTaskRepository) *MemoryTaskListRepository {
	return &MemoryTaskListRepository{
		taskLists: make(map[string]MemoryTaskList),
		taskRepo:  taskRepo, // Inicialize o campo taskRepo
	}
}

func (r *MemoryTaskListRepository) Create(taskList MemoryTaskList) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	taskListID := taskList.ID
	r.taskLists[taskListID] = taskList
	return taskListID, nil
}

func (r *MemoryTaskListRepository) GetByID(taskListID string) (*MemoryTaskList, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	taskList, exists := r.taskLists[taskListID]
	if !exists {
		return nil, errors.New("TaskList not found")
	}
	return &taskList, nil
}

func (r *MemoryTaskListRepository) Update(taskList MemoryTaskList) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.taskLists[taskList.ID]
	if !exists {
		return errors.New("TaskList not found")
	}

	r.taskLists[taskList.ID] = taskList
	return nil
}

func (r *MemoryTaskListRepository) Delete(taskListID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.taskLists[taskListID]
	if !exists {
		return errors.New("TaskList not found")
	}

	delete(r.taskLists, taskListID)
	return nil
}

func (r *MemoryTaskListRepository) AddTaskToList(taskID, taskListID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	taskList, exists := r.taskLists[taskListID]
	if !exists {
		return errors.New("TaskList not found")
	}

	taskList.Tasks = append(taskList.Tasks, taskID)
	r.taskLists[taskListID] = taskList
	return nil
}

func (r *MemoryTaskListRepository) GetTasksByList(taskListID string) ([]MemoryTask, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	taskList, exists := r.taskLists[taskListID]
	if !exists {
		return nil, errors.New("TaskList not found")
	}

	taskIDs := taskList.Tasks
	tasks, err := r.taskRepo.GetTasksByIDs(taskIDs) // Use o método GetTasksByIDs do taskRepo
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
