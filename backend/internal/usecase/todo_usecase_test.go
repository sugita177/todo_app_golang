package usecase

import (
	"context"
	"testing"
	"todo_app_golang/internal/domain"
)

// MockTodoRepository はテスト用の偽リポジトリ
type MockTodoRepository struct {
	CalledCreate bool
}

func (m *MockTodoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	m.CalledCreate = true
	return nil
}

func (m *MockTodoRepository) FetchAll(ctx context.Context) ([]*domain.Todo, error) {
	return nil, nil
}

func TestCreateTodo(t *testing.T) {
	mockRepo := &MockTodoRepository{}
	uc := NewTodoUseCase(mockRepo)

	t.Run("成功：タイトルがある場合", func(t *testing.T) {
		err := uc.CreateTodo(context.Background(), "買い物に行く")
		if err != nil {
			t.Errorf("エラーは発生しないはずですが発生しました: %v", err)
		}
		if !mockRepo.CalledCreate {
			t.Error("リポジトリのCreateメソッドが呼ばれていません")
		}
	})

	t.Run("失敗：タイトルが空の場合", func(t *testing.T) {
		err := uc.CreateTodo(context.Background(), "")
		if err != domain.ErrTitleEmpty {
			t.Errorf("期待されるエラー: %v, 実際のエラー: %v", domain.ErrTitleEmpty, err)
		}
	})
}
