package infrastructure

import (
	"context"
	"os"
	"testing"
	"time"
	"todo_app_golang/internal/domain"
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
