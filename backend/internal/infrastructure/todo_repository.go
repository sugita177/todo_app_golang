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
	// $1~$5 を使用し、RETURNING で ID と時間情報を取得
	query := `
		INSERT INTO todos (title, description, is_completed, priority, due_date, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		todo.Title, todo.Description, todo.IsCompleted, todo.Priority, todo.DueDate, todo.CreatedAt,
	).Scan(&todo.ID)

	return err
}

func (r *postgresTodoRepository) FetchAll(ctx context.Context) ([]*domain.Todo, error) {
	query := `SELECT id, title, description, is_completed, priority, due_date, created_at FROM todos ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*domain.Todo
	for rows.Next() {
		t := &domain.Todo{}
		// Scanの順序をSELECTと合わせる
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.IsCompleted, &t.Priority, &t.DueDate, &t.CreatedAt)
		if err != nil {
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

func (r *postgresTodoRepository) GetByID(ctx context.Context, id int) (*domain.Todo, error) {
	t := &domain.Todo{}
	query := `SELECT id, title, description, is_completed, priority, due_date, created_at FROM todos WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.Title, &t.Description, &t.IsCompleted, &t.Priority, &t.DueDate, &t.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// DB固有のエラーをドメインエラーに変換して返す
			return nil, domain.ErrTodoNotFound
		}
		return nil, err
	}
	return t, nil
}

func (r *postgresTodoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	query := `
		UPDATE todos 
		SET title = $1, description = $2, is_completed = $3, priority = $4, due_date = $5 
		WHERE id = $6`

	result, err := r.db.ExecContext(ctx, query,
		todo.Title, todo.Description, todo.IsCompleted, todo.Priority, todo.DueDate, todo.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
