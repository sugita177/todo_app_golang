package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"todo_app_golang/internal/infrastructure"
	"todo_app_golang/internal/usecase"
)

func main() {
	// 1. DB接続の初期化
	db, err := infrastructure.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// 2. 依存注入 (DI)
	repo := infrastructure.NewTodoRepository(db)
	todoUseCase := usecase.NewTodoUseCase(repo)

	// 3. ルーティング
	mux := http.NewServeMux()

	// ヘルスチェック
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	// TODO作成 (POST /todos)
	mux.HandleFunc("POST /todos", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Title string `json:"title"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := todoUseCase.CreateTodo(r.Context(), req.Title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Todo Created")
	})

	// TODO一覧取得 (GET /todos)
	mux.HandleFunc("GET /todos", func(w http.ResponseWriter, r *http.Request) {
		todos, err := todoUseCase.GetAllTodos(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(todos)
	})

	// 4. サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started at :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
