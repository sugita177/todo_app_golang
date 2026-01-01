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
	UpdateTodoStatus(ctx context.Context, id int, isCompleted bool) error
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

func (h *TodoHandler) UpdateTodoStatusHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// URLからIDを取得 (/todos/{id})
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// リクエストボディから新しい状態を取得
	var input struct {
		IsCompleted bool `json:"is_completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// UseCase の呼び出し
	if err := h.useCase.UpdateTodoStatus(ctx, id, input.IsCompleted); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
