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

- **Lint**: golangci-lintでコード品質チェック（main/PR）
- **Test**: ユニットテスト + カバレッジレポート（main/PR）
- **Build**: バイナリのビルド確認（main/PR）

実行タイミング：
- `main`ブランチへのpush時
- Pull Request作成・更新時

## 配布方法

現在はソースコードからのビルドのみ対応：

```bash
git clone https://github.com/SotaYamaguchi/meeting-toolkit.git
cd meeting-toolkit
./install.sh
```

将来的な改善案：
- [ ] リリース自動化（GitHub Actions）
- [ ] バイナリの自動ビルド・配布
- [ ] Homebrewでのインストール対応

## 注意事項

- `config.json` は顧客情報を含むため`.gitignore`対象
- サンプル設定（`config.sample.json`）のみコミット
- ドキュメント内の例は汎用的な名前を使用
