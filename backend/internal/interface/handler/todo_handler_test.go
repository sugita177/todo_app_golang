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
