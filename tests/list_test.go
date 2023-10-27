package tests

import (
	"botasks/internal/list"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTaskList(t *testing.T) {
	assert := assert.New(t)

	name := "My Task List"
	newList := list.CreateTaskList(name)

	assert.NotEmpty(newList.ID)
	assert.Equal(newList.Name, name)
	assert.Empty(newList.Tasks)
}
