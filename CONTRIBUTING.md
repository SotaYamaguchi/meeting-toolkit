# 開発者向けガイド

## 開発環境のセットアップ

```bash
# golangci-lintのインストール（macOS）
brew install golangci-lint

# pre-commitのインストール（推奨）
brew install pre-commit

# pre-commitフックを有効化
pre-commit install
```

## ローカルでのビルドとテスト

```bash
# ビルド
make build

# Lint実行
make lint

# ユニットテスト実行
make test

# または手動で
go test -v ./...
golangci-lint run

# カバレッジ付きテスト
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 設定ファイルを作成（開発用）
cp config.sample.json config.json
vim config.json  # プロジェクト情報を編集

# 動作確認
./mtg list
./mtg prep -project your-project
./mtg memo -project your-project
./mtg mail init -project your-project -type prep
./mtg mail -project your-project -type prep
```

## ディレクトリ構成

```
.
├── README.md                    # ユーザー向けドキュメント
├── CONTRIBUTING.md              # このファイル
├── CLAUDE.md                    # Claude Code向けガイド
├── install.sh                   # インストールスクリプト
├── Makefile                     # ビルド・テスト設定
├── go.mod                       # Go module定義
├── .golangci.yml                # Linter設定
├── .pre-commit-config.yaml      # pre-commit設定
├── .gitignore                   # Git除外設定
├── config.sample.json           # 設定ファイルのサンプル
├── main.go                      # エントリーポイント
├── main_test.go                 # テストコード
├── cmd/                         # サブコマンド実装
│   ├── root.go                  # ヘルプ表示
│   ├── prep.go                  # prepコマンド
│   ├── memo.go                  # memoコマンド
│   ├── mail.go                  # mailコマンド
│   ├── list.go                  # listコマンド
│   └── completion.go            # 補完スクリプト
├── pkg/                         # ビジネスロジック
│   ├── config/                  # 設定管理
│   │   └── config.go
│   ├── file/                    # ファイル操作
│   │   └── operations.go
│   └── mail/                    # メールテンプレート
│       └── template.go
└── templates/                   # サンプルテンプレート
    ├── project-a-prep.txt
    ├── project-a-memo.txt
    └── project-b-prep.txt
```

### パッケージ構造の説明

- **`main.go`**: エントリーポイント。サブコマンドディスパッチのみ
- **`cmd/`**: 各サブコマンドの実装。フラグ解析とpkgの呼び出し
- **`pkg/config/`**: config.jsonの読み書き、プレフィックス解決
- **`pkg/file/`**: ファイルのリネーム・集約処理
- **`pkg/mail/`**: メールテンプレートの解析・フォーマット・作成

## CI/CD

GitHub Actionsで以下を自動実行：

1. **Lint**: golangci-lintでコード品質チェック
2. **Test**: ユニットテスト + カバレッジレポート
3. **Build**: バイナリのビルド確認

## リリース手順

1. バージョンタグを作成
2. GitHub Actionsが自動でビルド・テスト
3. リリースノートを作成

## コーディング規約

- 全ての公開関数にGoDocコメントを追加
- パッケージレベルのコメントを各パッケージに記載
- エラーメッセージは日本語（ユーザー向け）
- 標準ライブラリのみ使用（外部依存なし）

## 注意事項

- `config.json` は `.gitignore` に含まれています（顧客情報保護）
- サンプル設定（`config.sample.json`）のみコミットしてください
- README内の例も汎用的な名前を使用してください
- 新機能追加時は対応するテストも追加してください
