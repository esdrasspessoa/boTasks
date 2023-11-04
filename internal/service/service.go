package service

import (
	"botasks/internal/list"
	"botasks/internal/repository"
	"botasks/internal/task"
	"errors"
	"time"
)


// TaskListService fornece a lógica de negócios para listas de tarefas e tarefas.
type TaskListService struct {
	taskListRepo repository.TaskListRepository
	taskRepo     repository.TaskRepository
}

// NewTaskListService cria uma nova instância de TaskListService.
func NewTaskListService(taskListRepo repository.TaskListRepository, taskRepo repository.TaskRepository) *TaskListService {
	return &TaskListService{
		taskListRepo: taskListRepo,
		taskRepo:     taskRepo,
	}
}

// CreateTaskList cria uma nova lista de tarefas.
func (s *TaskListService) CreateTaskList(name string) (string, error) {
	taskList := list.CreateTaskList(name) // Use a função CreateTaskList do pacote list
	taskListID, err := s.taskListRepo.Create(taskList)
	if err != nil {
		return "", err
	}
	return taskListID, nil
}

// UpdateTaskList atualiza uma lista de tarefas existente.
func (s *TaskListService) UpdateTaskList(taskListID, newName string) error {
	taskList, err := s.taskListRepo.GetByID(taskListID)
	if err != nil {
		return err
	}
	taskList.UpdateTaskList(newName) // Use o método UpdateTaskList do pacote list
	return s.taskListRepo.Update(*taskList)
}

// DeleteTaskList exclui uma lista de tarefas pelo ID.
func (s *TaskListService) DeleteTaskList(taskListID string) error {
	return s.taskListRepo.Delete(taskListID)
}

// AddTask cria uma nova tarefa e a adiciona a uma lista de tarefas.
func (s *TaskListService) AddTask(taskListID, title, description string, deadline time.Time) (string, error) {
	newTask := task.NewTask(title, description, deadline)
	if newTask == nil {
		return "", errors.New("invalid task parameters")
	}
	taskID, err := s.taskRepo.Create(*newTask)
	if err != nil {
		return "", err
	}
	err = s.taskListRepo.AddTaskToList(taskID, taskListID)
	if err != nil {
		return "", err
	}
	return taskID, nil
}

// UpdateTask atualiza uma tarefa existente.
func (s *TaskListService) UpdateTask(taskID, title, description string, deadline time.Time) error {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return err
	}
	task.UpdateTask(title, description, deadline) // Use o método UpdateTask do pacote task
	return s.taskRepo.Update(*task)
}

// DeleteTask exclui uma tarefa pelo ID.
func (s *TaskListService) DeleteTask(taskID string) error {
	return s.taskRepo.Delete(taskID)
}

// GetTaskList recupera uma lista de tarefas pelo ID.
func (s *TaskListService) GetTaskList(taskListID string) (*list.TaskList, error) {
	return s.taskListRepo.GetByID(taskListID)
}

// GetTasksByTaskList recupera todas as tarefas associadas a uma lista de tarefas.
func (s *TaskListService) GetTasksByTaskList(taskListID string) ([]task.Task, error) {
	listTasks, err := s.taskListRepo.GetTasksByList(taskListID)
	if err != nil {
		return nil, err
	}

	// Converte []list.Task para []task.Task
	tasks := make([]task.Task, len(listTasks))
	for i, listTask := range listTasks {
		tasks[i] = listTask.Task // Acessa o campo Task da estrutura composta list.Task
	}

	return tasks, nil
}


