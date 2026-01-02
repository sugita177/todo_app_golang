# Go & React Todo Application

Go (Clean Architecture) と React (TypeScript) で構築したフルスタックな TODO アプリケーションです。
Docker を使用して、開発環境を簡単に構築できるよう設計されています。

## 🚀 技術スタック

### Backend
- **Language:** Go 1.22+
- **Architecture:** Clean Architecture
- **Framework:** Standard Library (net/http)
- **Database:** PostgreSQL
- **Testing:** go test (Standard library)

### Frontend
- **Language:** TypeScript
- **Framework:** React 19
- **Build Tool:** Vite
- **Testing:** Vitest / React Testing Library / user-event

### Infrastructure
- **Container:** Docker / Docker Compose

## 🏗 アーキテクチャ (Backend)

バックエンドは保守性とテスト性を高めるため、クリーンアーキテクチャの考え方を採用しています。

- **Domain:** エンティティとリポジトリのインターフェース定義
- **UseCase:** ビジネスロジックの実行
- **Interface (Handler):** HTTPリクエストの解釈とレスポンスの返却
- **Infrastructure:** データベース接続などの外部実装



## 🛠 セットアップ

### 前提条件
- Docker / Docker Compose がインストールされていること

### 起動方法
```bash
# プロジェクトの起動
docker compose up -d --build
```

## 🌐 アクセス

起動後、以下のURLから各サービスにアクセスできます。

| サービス | URL | 説明 |
| :--- | :--- | :--- |
| **Frontend** | [http://localhost:5173](http://localhost:5173) | React 開発用サーバー |
| **Backend API** | [http://localhost:8080](http://localhost:8080) | Go REST API エンドポイント |

## 🧪 テストの実行

### バックエンド (Go)
```bash
docker compose exec app go test ./...
```

### フロントエンド (React)
```bash
# frontendディレクトリへ移動してから実行
npm test
```

## 📝 今後の実装予定

- [x] TODOの削除機能 (DELETE API)
- [x] TODOの完了状態の切り替え機能 (PATCH/PUT API)
- [ ] Tailwind CSSによるUI/UXの改善
- [ ] GitHub ActionsによるCI/CDの構築