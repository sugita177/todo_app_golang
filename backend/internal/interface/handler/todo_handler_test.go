package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo_app_golang/internal/domain"
)

// モック用の構造体を定義
type mockTodoUseCase struct {
	// テスト中に「エラーを発生させたい」などの制御ができるようにフィールドを持たせることが多い
	err error
}

// Interface を満たすようにメソッドを実装
func (m *mockTodoUseCase) CreateTodo(ctx context.Context, title string) error {
	return m.err // 設定したエラーを返す（成功時は nil）
}

func (m *mockTodoUseCase) GetAllTodos(ctx context.Context) ([]*domain.Todo, error) {
	return nil, m.err
}

func TestTodoHandler_CreateTodoHandler_Mock(t *testing.T) {
	// モックを生成（エラーなしの成功パターン）
	mockUC := &mockTodoUseCase{err: nil}
	h := NewTodoHandler(mockUC)

	jsonBody := []byte(`{"title": "Mockテストタスク"}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()

	// 実行
	h.CreateTodoHandler(rr, req)

	// 検証
	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %v", rr.Code)
	}
}

func TestTodoHandler_CreateTodoHandler_Error(t *testing.T) {
	// あえてエラーを返すモックを作成
	mockUC := &mockTodoUseCase{err: context.DeadlineExceeded}
	h := NewTodoHandler(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer([]byte(`{"title":"test"}`)))
	rr := httptest.NewRecorder()

	h.CreateTodoHandler(rr, req)

	// UseCaseがエラーを返した時、ハンドラーが 500 を返すかテスト
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %v", rr.Code)
	}
}
