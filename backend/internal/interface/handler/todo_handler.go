package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"todo_app_golang/internal/domain"
)

// ハンドラーが必要とする機能をインターフェースとして定義
type TodoUseCaseInterface interface {
	CreateTodo(ctx context.Context, title string) error
	GetAllTodos(ctx context.Context) ([]*domain.Todo, error)
	DeleteTodo(ctx context.Context, id int) error
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

func (h *TodoHandler) DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	// リクエストから Context を取得する
	ctx := r.Context()
	idStr := r.PathValue("id") // URLパラメータ {id} を取得
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.useCase.DeleteTodo(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
