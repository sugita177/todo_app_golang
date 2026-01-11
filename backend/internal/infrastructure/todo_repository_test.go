package infrastructure

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"
	"todo_app_golang/internal/domain"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // Postgresドライバ
	"github.com/stretchr/testify/assert"
)

// グローバル変数として保持し、各テストで共有する
var testDB *sql.DB

func TestMain(m *testing.M) {
	// 1. 環境変数の取得
	dsn := os.Getenv("TEST_DB_SOURCE")
	if dsn == "" {
		log.Println("TEST_DB_SOURCE is not set, skipping infrastructure tests")
		return
	}

	// 2. マイグレーションの実行
	// /app/migrations は Docker環境内のパスに合わせて調整してください
	migration, err := migrate.New("file:///app/migrations", dsn)
	if err != nil {
		log.Fatalf("Could not create migrate instance: %v", err)
	}

	// テスト開始前にクリーンな状態にするため一度DownしてからUp
	migration.Down()
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Could not run up migrations: %v", err)
	}

	// 3. DB接続の確立
	testDB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Could not connect to test database: %v", err)
	}

	// 4. 全てのテストを実行
	code := m.Run()

	// 5. 後片付け
	testDB.Close()
	os.Exit(code)
}

// ヘルパー関数: テストごとにクリーンなリポジトリを提供する
func setupRepository(t *testing.T) domain.TodoRepository {
	// テストごとにデータを消去したい場合はここで DELETE する
	_, err := testDB.Exec("DELETE FROM todos")
	if err != nil {
		t.Fatalf("Failed to clean table: %v", err)
	}
	return NewTodoRepository(testDB)
}

func TestPostgresTodoRepository_Create(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	// 3. テストデータの準備
	todo := &domain.Todo{
		Title:       "テストタスク",
		Description: "テストタスク詳細",
		IsCompleted: false,
		Priority:    "high",
		DueDate:     &time.Time{},
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
		testDB.Exec("DELETE FROM todos WHERE id = $1", todo.ID)
	})
}

func TestTodoRepository_FetchAll(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	// 1. テスト前にデータをクリア（他のテストの影響を排除）
	_, err := testDB.Exec("DELETE FROM todos")
	assert.NoError(t, err)

	// 2. テストデータを2件入れる（エラーを必ずチェックする）
	// 必須カラム（created_at等）がある場合はそれも指定する
	_, err = testDB.Exec(`INSERT INTO todos (title, description, is_completed, priority, due_date, created_at) 
	    VALUES ($1, $2, $3, $4, $5, $6)`,
		"Task 1", "Desc 1", false, "high", time.Now(), time.Now(),
	)
	if err != nil {
		t.Fatalf("テストデータの作成に失敗しました: %v", err)
	}

	_, err = testDB.Exec(`INSERT INTO todos (title, description, is_completed, priority, due_date, created_at) 
	    VALUES ($1, $2, $3, $4, $5, $6)`,
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
	// created_at の降順で取得される
	assert.Equal(t, "Task 2", todos[0].Title)
	assert.Equal(t, "Task 1", todos[1].Title)
}

func TestTodoRepository_Delete(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	// テストデータの準備
	var id int
	err := testDB.QueryRow(`INSERT INTO todos (title, description, is_completed, priority, due_date, created_at)
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
	testDB.QueryRow("SELECT count(*) FROM todos WHERE id = $1", id).Scan(&count)
	assert.Equal(t, 0, count)
}

func TestTodoRepository_UpdateStatus(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	// 1. テストデータの準備（未完了のタスクを作成）
	var id int
	err := testDB.QueryRow(`INSERT INTO todos (title, description, is_completed, priority, due_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		"Update Test Task", "Desc Update", false, "high", time.Now(), time.Now()).Scan(&id)
	assert.NoError(t, err)

	// 2. 実行：未完了(false)から完了(true)に更新
	err = repo.UpdateStatus(ctx, id, true)
	assert.NoError(t, err)

	// 3. 検証：DBの値が true になっているか確認
	var isCompleted bool
	err = testDB.QueryRow("SELECT is_completed FROM todos WHERE id = $1", id).Scan(&isCompleted)
	assert.NoError(t, err)
	assert.True(t, isCompleted)

	// 4. 実行：再度 false に戻せるか確認
	err = repo.UpdateStatus(ctx, id, false)
	assert.NoError(t, err)
	testDB.QueryRow("SELECT is_completed FROM todos WHERE id = $1", id).Scan(&isCompleted)
	assert.False(t, isCompleted)
}

func TestTodoRepository_GetByID(t *testing.T) {
	repo := setupRepository(t)
	ctx := context.Background()

	// 1. テストデータの準備
	dueDate := time.Now().Add(24 * time.Hour).Truncate(time.Microsecond)
	var id int
	err := testDB.QueryRow(`
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
		assert.Equal(t, domain.ErrTodoNotFound, err)
	})
}
