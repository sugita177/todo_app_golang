package infrastructure

import (
	"context"
	"database/sql"
	"todo_app_golang/internal/domain"
)

type postgresTodoRepository struct {
	db *sql.DB
}

// NewTodoRepository は Postgres 版のリポジトリを生成します
func NewTodoRepository(db *sql.DB) domain.TodoRepository {
	return &postgresTodoRepository{db: db}
}

func (r *postgresTodoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	query := `INSERT INTO todos (title, is_completed, created_at) VALUES ($1, $2, $3) RETURNING id`

	// $1, $2, $3 に値を流し込み、生成された ID を取得する
	err := r.db.QueryRowContext(ctx, query, todo.Title, todo.IsCompleted, todo.CreatedAt).Scan(&todo.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *postgresTodoRepository) FetchAll(ctx context.Context) ([]*domain.Todo, error) {
	query := `SELECT id, title, is_completed, created_at FROM todos`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*domain.Todo
	for rows.Next() {
		t := &domain.Todo{}
		if err := rows.Scan(&t.ID, &t.Title, &t.IsCompleted, &t.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}
