package infrastructure

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq" // 直接使わないが、ドライバを登録するために必要
)

// NewDB は PostgreSQL への接続を初期化します
func NewDB() (*sql.DB, error) {
	// docker-compose.yml の environment で設定したキー名を指定
	// DSN : Data Source Name
	dsn := os.Getenv("DB_SOURCE")
	if dsn == "" {
		// ローカル開発などで環境変数が設定されていない場合のデフォルト
		dsn = "postgresql://user:password@localhost:5432/todo?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// 接続確認
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
