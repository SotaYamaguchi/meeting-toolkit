# docpack - 顧客プロジェクトMTG支援ツール

顧客プロジェクトの定例MTG前後で使用するファイル整理ツールです。

## 特徴

- **単一ツール**: `mtg prep` と `mtg memo` でMTG前後の処理を統一的に実行
- **設定共有**: 1つの設定ファイルで全プロジェクトを管理
- **シンプル**: サブコマンド形式で直感的に使える

## インストール

### ワンライナーインストール

```bash
./install.sh
```

または

```bash
make install
```

### 手動インストール

```bash
cd mtg
make install
```

インストールすると：
- バイナリが `~/bin/mtg` に配置されます
- 設定ファイルが `~/.config/mtg/config.json` に配置されます

**重要**: インストール後、以下を `~/.zshrc` または `~/.bashrc` に追加してください。

```bash
export PATH="$HOME/bin:$PATH"
```

そして設定を反映：

```bash
source ~/.zshrc  # または source ~/.bashrc
```

## 使い方

### プロジェクト一覧を表示

```bash
mtg list
```

### MTG前の資料準備

```bash
# プロジェクト名で指定
mtg prep -project customer-a-project

# 別のディレクトリを指定
mtg prep -project customer-a-project -dir /path/to/directory
```

結果: `CUSTOMER_A_PREFIX_資料送付_20260318/` フォルダに資料が集約されます

### MTG後の議事メモ整理

```bash
# プロジェクト名で指定
mtg memo -project customer-a-project

# 別のディレクトリを指定
mtg memo -project customer-a-project -dir /path/to/directory
```

結果: `CUSTOMER_A_PREFIX_資料送付_20260318_MTG後/` フォルダに議事メモが集約されます

### ヘルプ表示

```bash
mtg help
```

## 設定ファイル

`~/.config/mtg/config.json` でプロジェクト名とプレフィックスのマッピングを管理します。

### 初回セットアップ

インストール後、設定ファイルを編集してプロジェクト情報を追加してください：

```bash
# エディタで設定ファイルを開く
vim ~/.config/mtg/config.json
# または
code ~/.config/mtg/config.json
```

### 設定例

```json
{
  "projects": {
    "customer-a-project": "CUSTOMER_A_PREFIX",
    "customer-b-project": "CUSTOMER_B_PREFIX",
    "internal-project": "INTERNAL_PREFIX"
  }
}
```

- **キー**: プロジェクト名（コマンドで使用）
- **値**: ファイル名のプレフィックス

新しいプロジェクトを追加する場合は、このファイルを編集してください。

**注意**: このファイルには顧客情報が含まれるため、Gitにコミットしないでください（`.gitignore`に含まれています）。

## アンインストール

```bash
cd mtg
make uninstall
```

## 元のシェルスクリプトからの移行

このツールは、プロジェクトごとに個別のシェルスクリプトを使っていた環境を統合したものです。

### 移行前（シェルスクリプト）

```bash
# プロジェクトAの会議前
./project-a_before_mtg_script.sh

# プロジェクトAの会議後
./project-a_after_mtg_script.sh
```

### 移行後（mtgツール）

```bash
# プロジェクトAの会議前
mtg prep -project project-a

# プロジェクトAの会議後
mtg memo -project project-a
```

## ディレクトリ構成

```
.
├── README.md                    # このファイル
├── install.sh                   # インストールスクリプト
├── Makefile                     # ルートMakefile
├── .gitignore                   # Git除外設定
└── mtg/                         # メインツール
    ├── main.go                  # ソースコード
    ├── go.mod                   # Go module定義
    ├── Makefile                 # ビルド設定
    ├── README.md                # 詳細ドキュメント
    └── config.sample.json       # 設定ファイルのサンプル
```

## 開発

### ローカルでのビルドとテスト

```bash
# ビルド
cd mtg
make build

# 設定ファイルを作成（開発用）
cp config.sample.json config.json
vim config.json  # プロジェクト情報を編集

# テスト実行
./mtg list
./mtg prep -project your-project
./mtg memo -project your-project
```

### GitHub公開時の注意

- `mtg/config.json` は `.gitignore` に含まれています（顧客情報保護）
- サンプル設定（`config.sample.json`）のみコミットしてください
- README内の例も汎用的な名前を使用してください
