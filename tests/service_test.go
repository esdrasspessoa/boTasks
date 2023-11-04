package tests

import (
	"botasks/internal/list"
	"botasks/internal/service"
	"botasks/internal/task"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks para o repositório de TaskList e Task.
type MockTaskListRepo struct {
	mock.Mock
}

type MockTaskRepo struct {
	mock.Mock
}

func (m *MockTaskListRepo) Create(taskList list.TaskList) (string, error) {
	args := m.Called(taskList)
	return args.String(0), args.Error(1)
}

func (m *MockTaskListRepo) GetByID(taskListID string) (*list.TaskList, error) {
	args := m.Called(taskListID)
	return args.Get(0).(*list.TaskList), args.Error(1)
}

func (m *MockTaskListRepo) Update(taskList list.TaskList) error {
	args := m.Called(taskList)
	return args.Error(0)
}

func (m *MockTaskListRepo) Delete(taskListID string) error {
	args := m.Called(taskListID)
	return args.Error(0)
}

func (m *MockTaskListRepo) AddTaskToList(taskID, taskListID string) error {
	args := m.Called(taskID, taskListID)
	return args.Error(0)
}

func (m *MockTaskListRepo) GetTasksByList(taskListID string) ([]list.Task, error) {
	args := m.Called(taskListID)
	return args.Get(0).([]list.Task), args.Error(1)
}

//

func (m *MockTaskRepo) Create(t task.Task) (string, error) {
	args := m.Called(t)
	return args.String(0), args.Error(1)
}

func (m *MockTaskRepo) GetByID(taskID string) (*task.Task, error) {
	args := m.Called(taskID)
	return args.Get(0).(*task.Task), args.Error(1)
}

func (m *MockTaskRepo) Update(t task.Task) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockTaskRepo) Delete(taskID string) error {
	args := m.Called(taskID)
	return args.Error(0)
}

// TestCreateTaskList verifica se o serviço cria uma lista de tarefas corretamente.
func TestCreateTaskList(t *testing.T) {
	mockTaskListRepo := new(MockTaskListRepo)
	mockTaskRepo := new(MockTaskRepo)
	s := service.NewTaskListService(mockTaskListRepo, mockTaskRepo)

	// Configure o mock para esperar qualquer TaskList e retorne um ID fictício.
	mockTaskListRepo.On("Create", mock.AnythingOfType("list.TaskList")).Return("id1", nil)

	taskListID, err := s.CreateTaskList("Lista de Tarefas")

	// Verifique se nenhum erro foi retornado e se o ID retornado é o esperado.
	assert.NoError(t, err)
	assert.Equal(t, "id1", taskListID)

	// Verifique se as expectativas configuradas no mock foram satisfeitas.
	mockTaskListRepo.AssertExpectations(t)
}

// TestGetTaskListByID verifica se o serviço recupera a lista de tarefas corretamente dado um ID válido.
func TestGetTaskListByID(t *testing.T) {
	mockTaskListRepo := new(MockTaskListRepo)
	mockTaskRepo := new(MockTaskRepo)
	s := service.NewTaskListService(mockTaskListRepo, mockTaskRepo)

	expectedTaskList := &list.TaskList{
		ID:    "id1",
		Name:  "Lista de Tarefas",
		Tasks: []list.Task{},
	}

	// Configurar o mock para retornar a lista de tarefas esperada quando o ID correspondente é fornecido.
	mockTaskListRepo.On("GetByID", "id1").Return(expectedTaskList, nil)

	// Chamar o método do serviço.
	taskList, err := s.GetTaskList("id1")

	// Verificar se nenhum erro foi retornado e se a lista de tarefas retornada é a esperada.
	assert.NoError(t, err)
	assert.Equal(t, expectedTaskList, taskList)

	// Verificar se as expectativas configuradas no mock foram satisfeitas.
	mockTaskListRepo.AssertExpectations(t)
}

func TestUpdateTaskList(t *testing.T) {
	mockTaskListRepo := new(MockTaskListRepo)
	mockTaskRepo := new(MockTaskRepo)
	s := service.NewTaskListService(mockTaskListRepo, mockTaskRepo)

	// Lista de tarefas que será recuperada e depois "atualizada"
	existingTaskList := list.TaskList{
		ID:   "id1",
		Name: "Lista de Tarefas Original",
	}

	// Configuração para o método GetByID do repositório ser chamado e para simular um sucesso na operação
	mockTaskListRepo.On("GetByID", "id1").Return(&existingTaskList, nil)
	// Configuração para o método Update do repositório ser chamado e para simular um sucesso na operação
	mockTaskListRepo.On("Update", mock.AnythingOfType("list.TaskList")).Return(nil)

	// Execução do método UpdateTaskList do serviço
	err := s.UpdateTaskList("id1", "Lista de Tarefas Atualizada")

	// Verificações: se nenhum erro foi retornado e se as expectativas do mock foram atendidas
	assert.NoError(t, err)
	mockTaskListRepo.AssertExpectations(t)
}

func TestUpdateNonExistentTaskList(t *testing.T) {
	mockTaskListRepo := new(MockTaskListRepo)
	s := service.NewTaskListService(mockTaskListRepo, nil) // taskRepo não é necessário para este teste

	// Configuração para o método GetByID do repositório ser chamado e para simular uma lista de tarefas inexistente
	mockTaskListRepo.On("GetByID", "nonexistent-id").Return((*list.TaskList)(nil), errors.New("TaskList not found"))

	// Execução do método UpdateTaskList do serviço
	err := s.UpdateTaskList("nonexistent-id", "New Name")

	// Verificações: se um erro foi retornado e se o erro é o esperado
	assert.Error(t, err)
	assert.EqualError(t, err, "TaskList not found")
	mockTaskListRepo.AssertExpectations(t)
}

func TestAddTask_Success(t *testing.T) {
	mockTaskListRepo := new(MockTaskListRepo)
	mockTaskRepo := new(MockTaskRepo)
	s := service.NewTaskListService(mockTaskListRepo, mockTaskRepo)

	taskListID := "existing-id"
	title := "New Task"
	description := "Task description"
	deadline := time.Now().Add(24 * time.Hour)

	// Configure os mocks para simular o comportamento esperado dos repositórios
	mockTaskRepo.On("Create", mock.AnythingOfType("task.Task")).Return("new-task-id", nil)
	mockTaskListRepo.On("AddTaskToList", "new-task-id", taskListID).Return(nil)

	// Executa o método AddTask do serviço
	taskID, err := s.AddTask(taskListID, title, description, deadline)

	// Verifica se nenhum erro foi retornado e se um ID de tarefa foi retornado
	assert.NoError(t, err)
	assert.Equal(t, "new-task-id", taskID)
	mockTaskRepo.AssertExpectations(t)
	mockTaskListRepo.AssertExpectations(t)
}

func TestAddTask_NonExistentTaskList(t *testing.T) {
	mockTaskListRepo := new(MockTaskListRepo)
	mockTaskRepo := new(MockTaskRepo)
	s := service.NewTaskListService(mockTaskListRepo, mockTaskRepo)

	title := "New Task"
	description := "New Task Description"
	deadline := time.Now().Add(24 * time.Hour) // Por exemplo, a deadline é daqui a 24 horas.
	taskListID := "nonexistent-id"

	// Aqui você precisa configurar o mock para o método Create, que é chamado dentro de AddTask.
	mockTaskRepo.On("Create", mock.AnythingOfType("task.Task")).Return("new-task-id", nil)
	// Configuração do mock para simular a tentativa de adicionar a uma lista inexistente.
	mockTaskListRepo.On("AddTaskToList", "new-task-id", taskListID).Return(errors.New("TaskList not found"))

	// Execução do método AddTask do serviço
	_, err := s.AddTask(taskListID, title, description, deadline)

	// Verificações: se um erro foi retornado e se o erro é o esperado.
	assert.Error(t, err)
	assert.EqualError(t, err, "TaskList not found")
	mockTaskListRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestAddTask_RepoError(t *testing.T) {
	mockTaskListRepo := new(MockTaskListRepo)
	mockTaskRepo := new(MockTaskRepo)
	s := service.NewTaskListService(mockTaskListRepo, mockTaskRepo)

	title := "New Task"
	description := "New Task Description"
	deadline := time.Now().Add(24 * time.Hour) // Por exemplo, a deadline é daqui a 24 horas.
	taskListID := "existing-id"

	// Configuração do mock para simular a criação de uma nova tarefa.
	mockTaskRepo.On("Create", mock.AnythingOfType("task.Task")).Return("new-task-id", nil)
	// Configuração do mock para simular um erro do repositório ao adicionar a tarefa à lista.
	mockTaskListRepo.On("AddTaskToList", "new-task-id", taskListID).Return(errors.New("internal error"))

	// Execução do método AddTask do serviço
	_, err := s.AddTask(taskListID, title, description, deadline)

	// Verificações: se um erro interno foi retornado.
	assert.Error(t, err)
	assert.EqualError(t, err, "internal error")
	mockTaskListRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestDeleteTaskList_Success(t *testing.T) {
	mockTaskListRepo := new(MockTaskListRepo)
	s := service.NewTaskListService(mockTaskListRepo, nil) // taskRepo não é necessário para este teste

	taskListID := "existing-id"

	// Configure o mock para simular a exclusão bem-sucedida da lista de tarefas
	mockTaskListRepo.On("Delete", taskListID).Return(nil)

	// Execução do método DeleteTaskList do serviço
	err := s.DeleteTaskList(taskListID)

	// Verificações: se nenhum erro foi retornado
	assert.NoError(t, err)
	mockTaskListRepo.AssertExpectations(t)
}

func TestDeleteTaskList_NonExistent(t *testing.T) {
	mockTaskListRepo := new(MockTaskListRepo)
	s := service.NewTaskListService(mockTaskListRepo, nil)

	taskListID := "nonexistent-id"

	// Configure o mock para simular que a lista de tarefas não existe
	mockTaskListRepo.On("Delete", taskListID).Return(errors.New("TaskList not found"))

	// Execução do método DeleteTaskList do serviço
	err := s.DeleteTaskList(taskListID)

	// Verificações: se um erro foi retornado e se o erro é o esperado
	assert.Error(t, err)
	assert.EqualError(t, err, "TaskList not found")
	mockTaskListRepo.AssertExpectations(t)
}

func TestDeleteTaskList_RepoError(t *testing.T) {
	mockTaskListRepo := new(MockTaskListRepo)
	s := service.NewTaskListService(mockTaskListRepo, nil)

	taskListID := "existing-id"

	// Configure o mock para simular um erro de repositório
	mockTaskListRepo.On("Delete", taskListID).Return(errors.New("internal error"))

	// Execução do método DeleteTaskList do serviço
	err := s.DeleteTaskList(taskListID)

	// Verificações: se um erro interno foi retornado
	assert.Error(t, err)
	assert.EqualError(t, err, "internal error")
	mockTaskListRepo.AssertExpectations(t)
}

func TestGetTasksByTaskList_Success(t *testing.T) {
	mockTaskListRepo := new(MockTaskListRepo)
	// mockTaskRepo não é necessário para este teste porque não interagimos com ele diretamente
	s := service.NewTaskListService(mockTaskListRepo, nil)

	taskListID := "existing-id"
	listTasks := []list.Task{
		{Task: task.Task{ID: "task1", Title: "Task 1", Description: "Description 1"}},
		{Task: task.Task{ID: "task2", Title: "Task 2", Description: "Description 2"}},
	}

	// Configure o mock para simular a recuperação bem-sucedida das tarefas como []list.Task
	mockTaskListRepo.On("GetTasksByList", taskListID).Return(listTasks, nil)

	// Execução do método GetTasksByTaskList do serviço
	returnedTasks, err := s.GetTasksByTaskList(taskListID)

	// Construir o resultado esperado como []task.Task
	expectedTasks := make([]task.Task, len(listTasks))
	for i, lt := range listTasks {
		expectedTasks[i] = lt.Task
	}

	// Verificações: se nenhum erro foi retornado e se as tarefas são as esperadas
	assert.NoError(t, err)
	assert.Equal(t, expectedTasks, returnedTasks)
	mockTaskListRepo.AssertExpectations(t)
}

func TestGetTasksByTaskList_NonExistent(t *testing.T) {
	mockTaskListRepo := new(MockTaskListRepo)
	s := service.NewTaskListService(mockTaskListRepo, nil)

	taskListID := "nonexistent-id"

	// Configure o mock para simular que a lista de tarefas não existe.
	// Certifique-se de que o primeiro retorno seja um slice de 'list.Task' que é explicitamente nil.
	mockTaskListRepo.On("GetTasksByList", taskListID).Return(([]list.Task)(nil), errors.New("TaskList not found"))

	// Execução do método GetTasksByTaskList do serviço.
	returnedTasks, err := s.GetTasksByTaskList(taskListID)

	// Verificações: se um erro foi retornado e se nenhuma tarefa foi retornada.
	assert.Error(t, err)
	assert.Nil(t, returnedTasks)
	assert.EqualError(t, err, "TaskList not found")
	mockTaskListRepo.AssertExpectations(t)
}

func TestGetTasksByTaskList_RepoError(t *testing.T) {
	mockTaskListRepo := new(MockTaskListRepo)
	s := service.NewTaskListService(mockTaskListRepo, nil)

	taskListID := "existing-id"

	// Configure o mock para simular um erro de repositório
	// e para retornar um slice 'nil' do tipo correto ([]list.Task).
	mockTaskListRepo.On("GetTasksByList", taskListID).Return(([]list.Task)(nil), errors.New("internal error"))

	// Execução do método GetTasksByTaskList do serviço
	returnedTasks, err := s.GetTasksByTaskList(taskListID)

	// Verificações: se um erro interno foi retornado e se nenhuma tarefa foi retornada
	assert.Error(t, err)
	assert.Nil(t, returnedTasks)
	assert.EqualError(t, err, "internal error")
	mockTaskListRepo.AssertExpectations(t)
}
