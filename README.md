# Webhook Server

Go言語とEchoフレームワークを使用したリアルタイムWebhookサーバーです。Server-Sent Events (SSE) を使用して、webhookで受信したメッセージをブラウザにリアルタイムで表示します。

## 機能

- ブラウザでメッセージ待機状態を表示
- POST `/webhook` でメッセージを受信
- Server-Sent Events (SSE) でリアルタイム更新
- ローカル開発とVercelデプロイに対応
- `go:embed`によるテンプレートファイル埋め込み

## 技術スタック

- **言語**: Go 1.22
- **フレームワーク**: Echo v4
- **デプロイ**: Vercel
- **リアルタイム通信**: Server-Sent Events (SSE)
- **テンプレート**: html/template + go:embed

## プロジェクト構成

```
demo/
├── api/
│   └── handler.go        # Vercel Function エントリーポイント (package handler)
├── pkg/
│   ├── handlers.go       # 共通ハンドラー関数 (package pkg)
│   ├── template.go       # テンプレートエンジン設定 + go:embed (package pkg)
│   └── templates/
│       └── index.html    # HTMLテンプレート (go:embedで埋め込み)
├── main.go               # ローカル実行用メインファイル (package main)
├── go.mod                # Go モジュール設定
├── vercel.json           # Vercel設定ファイル
└── README.md             # このファイル
```

## ローカル環境での実行

### 1. 依存関係のインストール

```bash
go mod tidy
```

### 2. ローカルサーバーの起動

`main.go` を直接実行してローカルサーバーを起動します：

```bash
# ローカルサーバー起動
go run main.go
```

### 3. ローカル環境でのテスト

#### ブラウザでの確認
```
http://localhost:3000
```

#### webhookのテスト

別のターミナルでcurlコマンドを実行：

```bash
# メッセージを送信
curl -X POST http://localhost:3000/webhook -d "message=Hello World"

# JSONフォーマットでの送信
curl -X POST http://localhost:3000/webhook \
  -H "Content-Type: application/json" \
  -d '{"message":"こんにちは世界"}'

# フォームデータでの送信
curl -X POST http://localhost:3000/webhook \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "message=テストメッセージ"
```

## Vercel環境でのデプロイとテスト

### 1. Vercelへのデプロイ

```bash
# Vercel CLIのインストール（初回のみ）
npm i -g vercel

# デプロイ
vercel

# または、GitHubリポジトリを接続してVercelで自動デプロイ
```

### 2. Vercel環境でのテスト

デプロイが完了すると、VercelからURLが提供されます（例：`https://your-project.vercel.app`）

#### ブラウザでの確認
```
https://your-project.vercel.app
```

#### webhookのテスト

```bash
# デプロイされたアプリのURLに置き換えてください
VERCEL_URL="https://your-project.vercel.app"

# メッセージを送信
curl -X POST $VERCEL_URL/webhook -d "message=Hello from Vercel"

# JSONフォーマットでの送信
curl -X POST $VERCEL_URL/webhook \
  -H "Content-Type: application/json" \
  -d '{"message":"Vercelからこんにちは"}'

# 日本語メッセージのテスト
curl -X POST $VERCEL_URL/webhook \
  -H "Content-Type: application/json" \
  -d '{"message":"リアルタイム更新テスト"}'
```

## 動作確認方法

1. **ブラウザでアクセス**
   - 初期状態では「メッセージ待機中」が表示される

2. **webhookでメッセージ送信**
   - curlコマンドでPOST リクエストを送信

3. **リアルタイム更新確認**
   - ブラウザのページが自動で更新され、送信したメッセージが表示される
   - ページの更新ボタンを押す必要なし

## API エンドポイント

### GET `/`
- ブラウザ表示用のHTMLページ
- 現在のメッセージとSSE接続を含む
- `go:embed`で埋め込まれたHTMLテンプレートを使用

### GET `/message`
- Polling用APIエンドポイント
- 現在のメッセージをJSON形式で返す
- レスポンス: `{"message": "現在のメッセージ"}`

### GET `/events`
- Server-Sent Events エンドポイント
- リアルタイム通信用（主にローカル環境）

### POST `/webhook`
- メッセージ受信エンドポイント
- パラメータ: `message` (string)
- 対応形式: JSON, form-data
- SSEクライアントにリアルタイムでブロードキャスト

## トラブルシューティング

### ローカル環境
- ポート3000が使用中の場合は、他のポートを使用
- `go mod tidy` で依存関係を更新

### Vercel環境
- デプロイエラーの場合は、`vercel.json` の設定を確認
- ログは Vercel ダッシュボードで確認可能

## 技術的な特徴

### go:embed による静的ファイル埋め込み
- HTMLテンプレートをバイナリに埋め込み
- Vercelでのファイルパス問題を解決
- シングルバイナリでの配布が可能

### 改善されたSSE接続管理
- **接続エラー対応**: Exponential backoff による自動再接続
- **メッセージ保持**: 接続エラー時もメッセージ状態を維持
- **Keepalive機能**: 30秒間隔のサーバーサイドkeepalive
- **適切なクリーンアップ**: 切断されたクライアントの自動削除

### 環境別リアルタイム通信
- **ローカル環境**: Server-Sent Events (SSE) でリアルタイム更新
- **Vercel環境**: SSE + 安定した再接続機能
- 両環境で一貫したユーザー体験
- 接続エラー時の適切なフォールバック処理

### パッケージ構成
- `pkg`パッケージで共通コードを管理
- Vercelの`internal`パッケージ制限を回避
- ローカルとVercelで同じビジネスロジックを共有

## 開発メモ

- SSEは長時間接続のため、Vercelの実行時間制限に注意
- 本番環境では適切なCORS設定を推奨
- メッセージはメモリ内保存のため、サーバー再起動で初期化される
- `go:embed`を使用しているため、テンプレート変更時は再ビルドが必要
- SSE接続エラーは正常な動作の一部で、自動再接続により透明に処理される