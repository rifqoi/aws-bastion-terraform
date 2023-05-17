package main

import (
	"context"

	"github.com/google/uuid"
)

type TodoService interface {
	GetTasks(ctx context.Context) ([]Todo, error)
	AddTask(ctx context.Context, task string) (*Todo, error)
	UpdateTaskStatus(ctx context.Context, id uuid.UUID) error
	RemoveTask(ctx context.Context, id uuid.UUID) error
}

type todoService struct {
	todoRepo TodoRepository
}

func NewTodoService(todoRepo TodoRepository) TodoService {
	return &todoService{
		todoRepo: todoRepo,
	}
}

func (t *todoService) GetTasks(ctx context.Context) ([]Todo, error) {
	todos, err := t.GetTasks(ctx)

	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (t *todoService) AddTask(ctx context.Context, task string) (*Todo, error) {
	todo := NewTodo(task, false)
	err := t.todoRepo.AddTask(ctx, todo)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (t *todoService) UpdateTaskStatus(ctx context.Context, id uuid.UUID) error {
	todo, err := t.todoRepo.GetTaskByID(ctx, id)
	if err != nil {
		return err
	}

	todo.ToggleDone()

	err = t.todoRepo.UpdateTask(ctx, *todo)
	if err != nil {
		return err
	}

	return nil
}

func (t *todoService) RemoveTask(ctx context.Context, id uuid.UUID) error {
	err := t.todoRepo.RemoveTask(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
