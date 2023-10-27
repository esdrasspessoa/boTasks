package list

import (
	"botasks/internal/task"

	"github.com/rs/xid"
)

type Task struct {
	task.Task
}

type TaskList struct {
	ID    string
	Name  string
	Tasks []Task
}

// CreateTaskList cria uma nova lista de tarefas com o nome especificado
func CreateTaskList(name string) TaskList {
	newList := TaskList{
		ID:    xid.New().String(),
		Name:  name,
		Tasks: []Task{},
	}
	return newList
}

// UpdateTaskList atualiza o nome de uma lista de tarefas existente
func (list *TaskList) UpdateTaskList(newName string) {
	list.Name = newName
}

// DeleteTaskList exclui uma lista de tarefas
func DeleteTaskList(lists []TaskList, listID string) []TaskList {
	for i, list := range lists {
		if list.ID == listID {
			lists = append(lists[:i], lists[i+1:]...)
			break
		}
	}
	return lists
}
