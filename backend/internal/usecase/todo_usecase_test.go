package usecase

import (
	"context"
	"testing"
	"todo_app_golang/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTodoRepository はテスト用の偽リポジトリ
type MockTodoRepository struct {
	mock.Mock // これを埋め込むことで On, Called, AssertExpectations が使えるようになる
}

func (m *MockTodoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

func (m *MockTodoRepository) FetchAll(ctx context.Context) ([]*domain.Todo, error) {
	args := m.Called(ctx)
	// return の型に合わせてキャストが必要
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Todo), args.Error(1)
}

func (m *MockTodoRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTodoRepository) UpdateStatus(ctx context.Context, id int, isCompleted bool) error {
	args := m.Called(ctx, id, isCompleted)
	return args.Error(0)
}

func TestCreateTodo(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	uc := NewTodoUseCase(mockRepo)
	ctx := context.Background()

	t.Run("成功：タイトルがある場合", func(t *testing.T) {
		// モックの期待値を設定 (Anyはどんな引数でも許容する場合に使用)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Todo")).Return(nil)

		err := uc.CreateTodo(ctx, "買い物に行く")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("失敗：タイトルが空の場合", func(t *testing.T) {
		err := uc.CreateTodo(ctx, "")
		assert.Equal(t, domain.ErrTitleEmpty, err)
	})
}

func TestGetAllTodos(t *testing.T) {
	ctx := context.Background()

	t.Run("成功：タスク一覧が取得できること", func(t *testing.T) {
		mockRepo := new(MockTodoRepository)
		useCase := NewTodoUseCase(mockRepo)
		// テスト用データを作成
		mockTodos := []*domain.Todo{
			{ID: 1, Title: "タスク1", IsCompleted: false},
			{ID: 2, Title: "タスク2", IsCompleted: true},
		}

		// モックの設定: 引数 ctx で呼ばれたら、mockTodos と nil を返す
		mockRepo.On("FetchAll", ctx).Return(mockTodos, nil)

		// 実行
		todos, err := useCase.GetAllTodos(ctx)

		// 検証
		assert.NoError(t, err)
		assert.Len(t, todos, 2)                 // 件数が正しいか
		assert.Equal(t, "タスク1", todos[0].Title) // 中身が正しいか
		mockRepo.AssertExpectations(t)
	})

	t.Run("成功：データが0件の場合に空の配列が返ること", func(t *testing.T) {
		mockRepo := new(MockTodoRepository)
		useCase := NewTodoUseCase(mockRepo)
		mockRepo.On("FetchAll", ctx).Return([]*domain.Todo{}, nil)

		todos, err := useCase.GetAllTodos(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, todos) // nilではなく空の配列であることを期待
		assert.Len(t, todos, 0)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteTodo(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	useCase := NewTodoUseCase(mockRepo)
	ctx := context.Background()
	targetID := 1

	// 「Deleteが呼ばれたらnilを返す」と定義
	mockRepo.On("Delete", ctx, targetID).Return(nil)

	err := useCase.DeleteTodo(ctx, targetID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t) // 実際に呼ばれたかチェック
}

func TestUpdateTodoStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("成功：完了状態を更新できること", func(t *testing.T) {
		mockRepo := new(MockTodoRepository)
		useCase := NewTodoUseCase(mockRepo)
		targetID := 10
		nextStatus := true

		// 期待値設定
		mockRepo.On("UpdateStatus", ctx, targetID, nextStatus).Return(nil)

		// 実行
		err := useCase.UpdateTodoStatus(ctx, targetID, nextStatus)

		// 検証
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
