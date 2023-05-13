package main

type TodoService interface {
	AddTodo(task string) (*Todo, error)
	UpdateTodo(id int) (*Todo, error)
	RemoveTodo(id int) error
}

type todoService struct {
	todoRepo TodoRepository
}

func NewTodoService(todoRepo TodoRepository) TodoService {
	return &todoService{
		todoRepo: todoRepo,
	}
}

func (t *todoService) AddTodo(task string) (*Todo, error) {
	return nil, nil
}

func (t *todoService) UpdateTodo(id int) (*Todo, error) {
	return nil, nil
}

func (t *todoService) RemoveTodo(id int) error {
	return nil
}
