package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 標準ライブラリのマルチプレクサ（ルーター）
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", mux)
}
