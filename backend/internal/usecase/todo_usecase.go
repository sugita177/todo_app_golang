package usecase

import (
	"context"
	"todo_app_golang/internal/domain"
)

type TodoUseCase struct {
	repo domain.TodoRepository
}

func NewTodoUseCase(repo domain.TodoRepository) *TodoUseCase {
	return &TodoUseCase{repo: repo}
}

// CreateTodo はバリデーションを行ってから保存を依頼します
func (u *TodoUseCase) CreateTodo(ctx context.Context, title string) error {
	todo, err := domain.NewTodo(title)
	if err != nil {
		return err
	}
	return u.repo.Create(ctx, todo)
}

func (u *TodoUseCase) GetAllTodos(ctx context.Context) ([]*domain.Todo, error) {
	return u.repo.FetchAll(ctx)
}

func (u *TodoUseCase) DeleteTodo(ctx context.Context, id int) error {
	return u.repo.Delete(ctx, id)
}
