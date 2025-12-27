package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"todo_app_golang/internal/domain"
)

// ハンドラーが必要とする機能をインターフェースとして定義
type TodoUseCaseInterface interface {
	CreateTodo(ctx context.Context, title string) error
	GetAllTodos(ctx context.Context) ([]*domain.Todo, error)
}

type TodoHandler struct {
	// 構造体 (*usecase.TodoUseCase) ではなく Interface を持つ
	useCase TodoUseCaseInterface
}

// 引数は Interface にする
func NewTodoHandler(uc TodoUseCaseInterface) *TodoHandler {
	return &TodoHandler{useCase: uc}
}

// CreateTodoHandler: POST /todos
func (h *TodoHandler) CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "無効なリクエストボディです", http.StatusBadRequest)
		return
	}

	err := h.useCase.CreateTodo(r.Context(), req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Todo Created"))
}

// GetAllTodosHandler: GET /todos
func (h *TodoHandler) GetAllTodosHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := h.useCase.GetAllTodos(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}
