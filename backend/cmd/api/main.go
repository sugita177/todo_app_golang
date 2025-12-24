package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"todo_app_golang/internal/usecase"
	// 後ほど実装する infrastructure を読み込めるよう準備
	// "todo_app_golang/internal/infrastructure"
)

func main() {
	// 1. 本来はここでDB接続を初期化します
	// db, err := infrastructure.NewDB()
	// if err != nil { log.Fatal(err) }
	// repo := infrastructure.NewTodoRepository(db)

	// 現時点では、動作確認のために一旦 nil（またはモック）を渡しておきます
	// 後で本物のリポジトリに差し替えます
	todoUseCase := usecase.NewTodoUseCase(nil)

	// 2. 標準ライブラリ (Go 1.22+) のマルチプレクサ
	mux := http.NewServeMux()

	// ヘルスチェック用エンドポイント
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	// TODO作成エンドポイント (例)
	mux.HandleFunc("POST /todos", func(w http.ResponseWriter, r *http.Request) {
		// ここで本来はJSONをパースして usecase.CreateTodo を呼び出します
		fmt.Fprintln(w, "Todo Created (Stub)")
	})

	// 3. サーバーの起動設定
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Server started at %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
