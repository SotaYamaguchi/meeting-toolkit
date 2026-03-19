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
cd mtg
make build

# Lint実行
make lint

# ユニットテスト実行
make test

# または手動で
go test -v
golangci-lint run

# カバレッジ付きテスト
go test -v -race -coverprofile=coverage.out
go tool cover -html=coverage.out

# 設定ファイルを作成（開発用）
cp config.sample.json config.json
vim config.json  # プロジェクト情報を編集

# 動作確認
./mtg list
./mtg prep -project your-project
./mtg memo -project your-project
```

## ディレクトリ構成

```
.
├── README.md                    # ユーザー向けドキュメント
├── CONTRIBUTING.md              # このファイル
├── install.sh                   # インストールスクリプト
├── Makefile                     # ルートMakefile
├── .golangci.yml                # Linter設定
├── .pre-commit-config.yaml      # pre-commit設定
├── .gitignore                   # Git除外設定
└── mtg/                         # メインツール
    ├── main.go                  # ソースコード
    ├── main_test.go             # テストコード
    ├── go.mod                   # Go module定義
    ├── Makefile                 # ビルド設定
    ├── README.md                # 詳細ドキュメント
    └── config.sample.json       # 設定ファイルのサンプル
```

## CI/CD

GitHub Actionsで以下を自動実行：

1. **Lint**: golangci-lintでコード品質チェック
2. **Test**: ユニットテスト + カバレッジレポート
3. **Build**: バイナリのビルド確認

## リリース手順

1. バージョンタグを作成
2. GitHub Actionsが自動でビルド・テスト
3. リリースノートを作成

## 注意事項

- `mtg/config.json` は `.gitignore` に含まれています（顧客情報保護）
- サンプル設定（`config.sample.json`）のみコミットしてください
- README内の例も汎用的な名前を使用してください
