package domain

import (
	"context"
	"errors"
	"time"
)

// Todo はタスクを表すエンティティです
type Todo struct {
	ID          int        `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"` // 詳細説明用
	IsCompleted bool       `json:"is_completed" db:"is_completed"`
	Priority    string     `json:"priority" db:"priority"` // 'low', 'medium', 'high'
	DueDate     *time.Time `json:"due_date" db:"due_date"` // 期限（未設定を許容するためポインタ）
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"` // 更新日時も持っておくと便利です
}

// TodoRepository はデータ操作に関するインターフェースです
type TodoRepository interface {
	Create(ctx context.Context, todo *Todo) error
	FetchAll(ctx context.Context) ([]*Todo, error)
	Delete(ctx context.Context, id int) error
	UpdateStatus(ctx context.Context, id int, isCompleted bool) error
	GetByID(ctx context.Context, id int) (*Todo, error)
	Update(ctx context.Context, todo *Todo) error
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
