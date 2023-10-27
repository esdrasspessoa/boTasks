package tests

import (
	"botasks/internal/task"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	// Caso de teste válido
	t1 := task.NewTask("Task 1", "Description 1", time.Now().Add(time.Hour))
	assert.NotNil(t, t1)
	assert.NotEmpty(t, t1.ID)
	assert.Equal(t, "Task 1", t1.Title)
	assert.Equal(t, "Description 1", t1.Description)
	assert.True(t, t1.Deadline.After(time.Now()))

	// Caso de teste inválido: título em branco
	t2 := task.NewTask("", "Description 2", time.Now().Add(time.Hour))
	assert.Nil(t, t2)

	// Caso de teste inválido: prazo no passado
	t3 := task.NewTask("Task 3", "Description 3", time.Now().Add(-time.Hour))
	assert.Nil(t, t3)

	// Caso de teste inválido: título em branco e prazo no passado
	t4 := task.NewTask("", "Description 4", time.Now().Add(-time.Hour))
	assert.Nil(t, t4)
}
