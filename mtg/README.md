# docpack

顧客プロジェクトの定例MTG前後でファイルを整理するCLIツールです。

## 特徴

- **単一ツール**: `mtg prep` と `mtg memo` でMTG前後の処理を統一的に実行
- **設定共有**: 1つの設定ファイルで全プロジェクトを管理
- **シンプル**: サブコマンド形式で直感的に使える

## インストール

```bash
# ビルドしてインストール
make install
```

インストールすると：
- バイナリが `~/bin/mtg` に配置されます
- 設定ファイルが `~/.config/mtg/config.json` に配置されます
- zsh用のタブ補完スクリプトが `~/.zsh/completions/_mtg` に配置されます
- 既存の `mtg-prep`、`mtg-memo` を自動的に削除します

**注意**: `~/bin` がPATHに含まれていることを確認してください。

```bash
export PATH="$HOME/bin:$PATH"
```

**タブ補完を有効にする (zsh)**:

`~/.zshrc` に以下を追加してシェルを再起動してください：

```bash
fpath=(~/.zsh/completions $fpath)
autoload -Uz compinit && compinit
```

## 使い方

### プロジェクト一覧を表示

```bash
mtg list
```

### MTG前の資料準備

```bash
# プロジェクト名で指定
mtg prep -project project-a

# 別のディレクトリを指定
mtg prep -project project-a -dir /path/to/directory

# プレフィックスを直接指定
mtg prep -prefix PREFIX_A
```

処理内容:
1. ファイル名の "main" を日付に置換（例: `20260318`）
2. `{PREFIX}_資料送付_{日付}/` フォルダを作成（例: `PREFIX_A_資料送付_20260318/`）
3. ファイルを集約

### MTG後の議事メモ整理

```bash
# プロジェクト名で指定
mtg memo -project project-a

# 別のディレクトリを指定
mtg memo -project project-a -dir /path/to/directory

# プレフィックスを直接指定
mtg memo -prefix PREFIX_A
```

処理内容:
1. ファイル名の "main" を日付+_MTG後に置換（例: `20260318_MTG後`）
2. `{PREFIX}_資料送付_{日付}_MTG後/` フォルダを作成（例: `PREFIX_A_資料送付_20260318_MTG後/`）
3. ファイルを集約

## 設定ファイル

`~/.config/mtg/config.json` でプロジェクト名とプレフィックスのマッピングを管理します。

```json
{
  "projects": {
    "project-a": "PREFIX_A",
    "project-b": "PREFIX_B",
    "project-c": "PREFIX_C"
  }
}
```

**注意**: 実際の顧客情報を含む `config.json` は `.gitignore` で除外されます。`config.sample.json` をテンプレートとして使用してください。

新しいプロジェクトを追加する場合は、このファイルを編集してください。

## アンインストール

```bash
make uninstall
```

## 元のツールからの移行

以前の `mtg-prep` と `mtg-memo` から自動的に移行されます：

- `mtg-prep -project <name>` → `mtg prep -project <name>`
- `mtg-memo -project <name>` → `mtg memo -project <name>`
- `~/.config/mtg-prep/config.json` → `~/.config/mtg/config.json`

## 開発

```bash
# ビルド
make build

# 設定ファイルを作成（開発用）
cp config.sample.json config.json
vim config.json  # プロジェクト情報を編集

# テスト実行
./mtg list
./mtg prep -project your-project
./mtg memo -project your-project
```

## ヘルプ

```bash
mtg help
```
