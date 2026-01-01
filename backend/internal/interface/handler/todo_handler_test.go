package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo_app_golang/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// testify/mock スタイルに統一
type mockTodoUseCase struct {
	mock.Mock
}

func (m *mockTodoUseCase) CreateTodo(ctx context.Context, title string) error {
	args := m.Called(ctx, title)
	return args.Error(0)
}

func (m *mockTodoUseCase) GetAllTodos(ctx context.Context) ([]*domain.Todo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Todo), args.Error(1)
}

func (m *mockTodoUseCase) DeleteTodo(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockTodoUseCase) UpdateTodoStatus(ctx context.Context, id int, isCompleted bool) error {
	args := m.Called(ctx, id, isCompleted)
	return args.Error(0)
}

// --- テストケース ---

func TestTodoHandler_CreateTodoHandler_Mock(t *testing.T) {
	mockUC := new(mockTodoUseCase)
	h := NewTodoHandler(mockUC)

	// 設定: CreateTodo が呼ばれたら nil (成功) を返す
	mockUC.On("CreateTodo", mock.Anything, "Mockテストタスク").Return(nil)

	jsonBody := []byte(`{"title": "Mockテストタスク"}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()

	h.CreateTodoHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockUC.AssertExpectations(t)
}

func TestTodoHandler_CreateTodoHandler_Error(t *testing.T) {
	mockUC := new(mockTodoUseCase)
	h := NewTodoHandler(mockUC)

	// 設定: エラーを返すようにする
	mockUC.On("CreateTodo", mock.Anything, "test").Return(context.DeadlineExceeded)

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer([]byte(`{"title":"test"}`)))
	rr := httptest.NewRecorder()

	h.CreateTodoHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestTodoHandler_DeleteTodoHandler(t *testing.T) {
	mockUC := new(mockTodoUseCase)
	handler := NewTodoHandler(mockUC)

	// ID:1 の削除リクエスト
	req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	req.SetPathValue("id", "1") // Go 1.22+ パラメータ
	w := httptest.NewRecorder()

	// 期待値設定: ID 1 が渡されたら成功を返す
	mockUC.On("DeleteTodo", mock.Anything, 1).Return(nil)

	handler.DeleteTodoHandler(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockUC.AssertExpectations(t)
}

func TestTodoHandler_UpdateTodoStatusHandler(t *testing.T) {
	t.Run("成功：ステータスを更新できること", func(t *testing.T) {
		mockUC := new(mockTodoUseCase)
		h := NewTodoHandler(mockUC)

		// テストデータ
		targetID := 5
		nextStatus := true

		// ボディの作成
		jsonBody := []byte(`{"is_completed": true}`)
		req := httptest.NewRequest(http.MethodPatch, "/todos/5", bytes.NewBuffer(jsonBody))

		// Go 1.22+ のパスパラメータをシミュレート
		req.SetPathValue("id", "5")

		rr := httptest.NewRecorder()

		// 期待値設定: ID=5, Status=true で呼ばれることを期待
		mockUC.On("UpdateTodoStatus", mock.Anything, targetID, nextStatus).Return(nil)

		// 実行
		h.UpdateTodoStatusHandler(rr, req)

		// 検証
		assert.Equal(t, http.StatusNoContent, rr.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("失敗：不正なJSONボディの場合に400を返すこと", func(t *testing.T) {
		mockUC := new(mockTodoUseCase)
		h := NewTodoHandler(mockUC)

		// 壊れたJSONを送る
		req := httptest.NewRequest(http.MethodPatch, "/todos/5", bytes.NewBuffer([]byte(`{invalid-json}`)))
		req.SetPathValue("id", "5")
		rr := httptest.NewRecorder()

		// 実行
		h.UpdateTodoStatusHandler(rr, req)

		// 検証
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
