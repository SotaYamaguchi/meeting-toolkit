# docpack

[![Test](https://github.com/SotaYamaguchi/docpack/actions/workflows/test.yml/badge.svg)](https://github.com/SotaYamaguchi/docpack/actions/workflows/test.yml)

MTG前後の資料ファイルを自動で整理するCLIツール。

## インストール

```bash
./install.sh
```

自動セットアップを選択すると、シェルを再起動するだけで使えます：

```bash
exec zsh
```

## 使い方

```bash
# プロジェクト一覧
mtg list

# MTG前の資料準備
mtg prep -project your-project

# MTG後の議事メモ整理
mtg memo -project your-project
```

## 初期設定

`~/.config/mtg/config.json` を編集してプロジェクトを追加：

```json
{
  "projects": {
    "project-a": "PREFIX_A",
    "project-b": "PREFIX_B"
  }
}
```

```bash
vim ~/.config/mtg/config.json
```

## 詳細

- タブ補完対応（zsh）
- ヘルプ: `mtg help`
- アンインストール: `cd mtg && make uninstall`
- 開発者向け: [CONTRIBUTING.md](CONTRIBUTING.md)
