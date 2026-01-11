package infrastructure

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"
	"todo_app_golang/internal/domain"

	_ "github.com/lib/pq" // Postgresドライバ
	"github.com/stretchr/testify/assert"
)

func TestPostgresTodoRepository_Create(t *testing.T) {
	// 開発用(DB_SOURCE)ではなく、テスト用(TEST_DB_SOURCE)を取得
	dsn := os.Getenv("TEST_DB_SOURCE")
	if dsn == "" {
		t.Skip("TEST_DB_SOURCE が設定されていないためテストをスキップします")
	}

	// db.go 内の sql.Open を、このテスト用DSNで実行するようにする
	// ※ NewDB() が os.Getenv("DB_SOURCE") を直接見ている場合は、
	// 以下のように一時的に環境変数を上書きするか、NewDBを引数付きにするのが一般的です。

	originalDSN := os.Getenv("DB_SOURCE")
	os.Setenv("DB_SOURCE", dsn)               // 一時的に差し替え
	defer os.Setenv("DB_SOURCE", originalDSN) // テスト後に戻す

	db, err := NewDB()
	if err != nil {
		t.Fatalf("テスト用DB接続失敗: %v", err)
	}
	defer db.Close()

	// テスト開始前にテーブルを確実に用意する（存在しなければ作成）
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS todos (
	    id SERIAL PRIMARY KEY, 
	    title TEXT NOT NULL, 
	    description TEXT,              -- 追加
	    is_completed BOOLEAN NOT NULL, 
	    priority VARCHAR(10),          -- 追加
	    due_date TIMESTAMP WITH TIME ZONE, -- 追加
	    created_at TIMESTAMP WITH TIME ZONE NOT NULL
	);`)
	if err != nil {
		t.Fatalf("テーブル作成失敗: %v", err)
	}

	repo := NewTodoRepository(db)
	ctx := context.Background()

	// 3. テストデータの準備
	todo := &domain.Todo{
		Title:       "テストタスク",
		IsCompleted: false,
		CreatedAt:   time.Now(),
	}

	// 4. 実行
	t.Run("Todoが正常に保存されIDが採番されること", func(t *testing.T) {
		err := repo.Create(ctx, todo)
		if err != nil {
			t.Fatalf("Create失敗: %v", err)
		}

		if todo.ID == 0 {
			t.Error("IDが採番されていません")
		}
	})

	// 5. 後片付け（重要！）
	// テストで作ったデータを削除して、DBを元の状態に戻す
	t.Cleanup(func() {
		db.Exec("DELETE FROM todos WHERE id = $1", todo.ID)
	})
}

func TestTodoRepository_FetchAll(t *testing.T) {
	dsn := os.Getenv("TEST_DB_SOURCE")
	if dsn == "" {
		dsn = "host=db_test port=5432 user=test_user password=test_password dbname=todo_test sslmode=disable"
	}
	db, _ := sql.Open("postgres", dsn)
	defer db.Close()

	repo := NewTodoRepository(db)
	ctx := context.Background()

	// 1. テスト前にデータをクリア（他のテストの影響を排除）
	_, err := db.Exec("DELETE FROM todos")
	assert.NoError(t, err)

	// 2. テストデータを2件入れる（エラーを必ずチェックする）
	// 必須カラム（created_at等）がある場合はそれも指定する
	_, err = db.Exec(`INSERT INTO todos (title, description, is_completed, priority, due_date, created_at) 
	    VALUES ($1, $2, $3, $4, $5, $6)`,
		"Task 1", "Desc 1", false, "high", time.Now(), time.Now(),
		"Task 2", "Desc 2", true, "low", time.Now(), time.Now(),
	)
	if err != nil {
		t.Fatalf("テストデータの作成に失敗しました: %v", err)
	}

	// 3. 実行
	todos, err := repo.FetchAll(ctx)

	// 4. 検証
	assert.NoError(t, err)
	// 具体的な件数でチェックするのがベストです
	assert.Equal(t, 2, len(todos), "取得されたタスク数が一致しません")
	assert.Equal(t, "Task 1", todos[0].Title)
}

func TestTodoRepository_Delete(t *testing.T) {
	// DSNを環境変数から取得するようにしておくと安全です
	dsn := os.Getenv("TEST_DB_SOURCE")
	if dsn == "" {
		dsn = "host=db_test port=5432 user=test_user password=test_password dbname=todo_test sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("DB接続失敗: %v", err)
	}
	defer db.Close()

	repo := NewTodoRepository(db)
	ctx := context.Background()

	// テストデータの準備
	var id int
	err = db.QueryRow(`INSERT INTO todos (title, description, is_completed, priority, due_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		"Test Delete", "Desc Delete", false, "high", time.Now(), time.Now()).Scan(&id)
	if err != nil {
		t.Fatalf("テストデータ作成失敗: %v", err)
	}

	// 実行
	err = repo.Delete(ctx, id)

	// 検証
	assert.NoError(t, err)

	// DBから消えているか確認
	var count int
	db.QueryRow("SELECT count(*) FROM todos WHERE id = $1", id).Scan(&count)
	assert.Equal(t, 0, count)
}

func TestTodoRepository_UpdateStatus(t *testing.T) {
	dsn := os.Getenv("TEST_DB_SOURCE")
	if dsn == "" {
		dsn = "host=db port=5432 user=user password=password dbname=todo sslmode=disable"
	}
	db, _ := sql.Open("postgres", dsn)
	defer db.Close()

	repo := NewTodoRepository(db)
	ctx := context.Background()

	// 1. テストデータの準備（未完了のタスクを作成）
	var id int
	err := db.QueryRow(`INSERT INTO todos (title, description, is_completed, priority, due_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		"Update Test Task", "Desc Update", false, "high", time.Now(), time.Now()).Scan(&id)
	assert.NoError(t, err)

	// 2. 実行：未完了(false)から完了(true)に更新
	err = repo.UpdateStatus(ctx, id, true)
	assert.NoError(t, err)

	// 3. 検証：DBの値が true になっているか確認
	var isCompleted bool
	err = db.QueryRow("SELECT is_completed FROM todos WHERE id = $1", id).Scan(&isCompleted)
	assert.NoError(t, err)
	assert.True(t, isCompleted)

	// 4. 実行：再度 false に戻せるか確認
	err = repo.UpdateStatus(ctx, id, false)
	assert.NoError(t, err)
	db.QueryRow("SELECT is_completed FROM todos WHERE id = $1", id).Scan(&isCompleted)
	assert.False(t, isCompleted)
}

func TestTodoRepository_GetByID(t *testing.T) {
	dsn := os.Getenv("TEST_DB_SOURCE")
	db, _ := sql.Open("postgres", dsn)
	defer db.Close()

	repo := NewTodoRepository(db)
	ctx := context.Background()

	// 1. テストデータの準備
	dueDate := time.Now().Add(24 * time.Hour).Truncate(time.Microsecond)
	var id int
	err := db.QueryRow(`
        INSERT INTO todos (title, description, is_completed, priority, due_date, created_at) 
        VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		"Detail Test", "Description here", true, "low", dueDate, time.Now(),
	).Scan(&id)
	assert.NoError(t, err)

	t.Run("存在するIDを指定した場合、全てのフィールドが取得できること", func(t *testing.T) {
		todo, err := repo.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, id, todo.ID)
		assert.Equal(t, "Detail Test", todo.Title)
		assert.Equal(t, "Description here", todo.Description)
		assert.Equal(t, "low", todo.Priority)
		assert.True(t, todo.IsCompleted)
		// Timeの比較は .Equal を使用（タイムゾーンの差異を許容）
		assert.True(t, dueDate.Equal(*todo.DueDate))
	})

	t.Run("存在しないIDを指定した場合、エラーが返ること", func(t *testing.T) {
		_, err := repo.GetByID(ctx, 99999)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})
}
