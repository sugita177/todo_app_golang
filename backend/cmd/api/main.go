package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	// 4. サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started at :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
