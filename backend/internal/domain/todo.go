package domain

import (
	"context"
	"errors"
	"time"
)

// Todo はタスクを表すエンティティです
type Todo struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	IsCompleted bool      `json:"is_completed" db:"is_completed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// TodoRepository はデータ操作に関するインターフェースです
type TodoRepository interface {
	Create(ctx context.Context, todo *Todo) error
	FetchAll(ctx context.Context) ([]*Todo, error)
	Delete(ctx context.Context, id int) error
	UpdateStatus(ctx context.Context, id int, isCompleted bool) error
}

var (
	ErrTitleEmpty = errors.New("タイトルを入力してください")
)

// NewTodo は新しいTodoを生成する際のビジネスルールを適用します
func NewTodo(title string) (*Todo, error) {
	if title == "" {
		return nil, ErrTitleEmpty
	}
	return &Todo{
		Title:       title,
		IsCompleted: false,
		CreatedAt:   time.Now(),
	}, nil
}
