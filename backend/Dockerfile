FROM golang:1.25-alpine

WORKDIR /app

# ホットリロード用にAirをインストール（開発効率が上がります）
RUN go install github.com/air-verse/air@latest

COPY go.mod ./
# 今回は標準ライブラリのみなので go.sum は不要な場合もあります
RUN go mod download

COPY . .

# Airで実行
CMD ["air", "-c", ".air.toml"]