package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"

	"todo_app_golang/internal/infrastructure"
	"todo_app_golang/internal/interface/handler"
	"todo_app_golang/internal/usecase"
)

func main() {
	// 1. DB接続 (インフラ層)
	db, err := infrastructure.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// 2. 依存注入 (DI)
	repo := infrastructure.NewTodoRepository(db)
	todoUseCase := usecase.NewTodoUseCase(repo)
	todoHandler := handler.NewTodoHandler(todoUseCase) // ハンドラーを生成

	// 3. ルーティング
	mux := http.NewServeMux()

	// インターフェース層のメソッドを紐付け
	mux.HandleFunc("POST /todos", todoHandler.CreateTodoHandler)
	mux.HandleFunc("GET /todos", todoHandler.GetAllTodosHandler)

	// CORS 設定
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, // フロントエンドのURLを許可
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	// mux を cors ハンドラーで包む
	handler := c.Handler(mux)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
