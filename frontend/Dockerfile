FROM node:24-alpine

WORKDIR /app

# 先に package.json をコピーしてライブラリをインストール
COPY package*.json ./
RUN npm install

# 残りのファイルをコピー
COPY . .

# Viteのデフォルトポート
EXPOSE 5173

# 開発サーバーを起動（--host オプションが必要）
CMD ["npm", "run", "dev", "--", "--host"]