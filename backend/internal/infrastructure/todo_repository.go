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

func (r *postgresTodoRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM todos WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id) // Exec ではなく ExecContext を使うのがベスト
	return err
}

func (r *postgresTodoRepository) UpdateStatus(ctx context.Context, id int, isCompleted bool) error {
	query := `UPDATE todos SET is_completed = $1 WHERE id = $2`

	// ExecContext を使用してクエリを実行
	result, err := r.db.ExecContext(ctx, query, isCompleted, id)
	if err != nil {
		return err
	}

	// 念のため、更新された行数を確認（存在しないIDが指定された場合のケア）
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows // 1行も更新されなかったらエラーとする
	}

	return nil
}

func (r *postgresTodoRepository) GetByID(id int) (*domain.Todo, error) {
	var todo domain.Todo
	query := "SELECT id, title, description, is_completed, priority, due_date, created_at FROM todos WHERE id = ?"
	err := r.db.QueryRow(query, id).Scan(
		&todo.ID, &todo.Title, &todo.Description, &todo.IsCompleted,
		&todo.Priority, &todo.DueDate, &todo.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}
