#!/bin/bash
set -e

echo "mtg - 顧客プロジェクトMTG支援ツールをインストールします"
echo ""

# mtgディレクトリに移動してインストール
cd "$(dirname "$0")/mtg"
make install

echo ""
echo "=========================================="
echo "シェル設定を自動セットアップしますか？"
echo "=========================================="
echo ""
echo "以下の設定を ~/.zshrc に追加します:"
echo "  - PATH設定 (~/bin)"
echo "  - タブ補完設定"
echo ""
read -p "自動セットアップしますか？ (y/N): " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    ZSHRC="$HOME/.zshrc"

    # .zshrcが存在しない場合は作成
    if [ ! -f "$ZSHRC" ]; then
        touch "$ZSHRC"
        echo "~/.zshrc を作成しました"
    fi

    # PATH設定を追加
    if ! grep -q 'export PATH="$HOME/bin:$PATH"' "$ZSHRC"; then
        echo "" >> "$ZSHRC"
        echo "# mtg tool - PATH設定" >> "$ZSHRC"
        echo 'export PATH="$HOME/bin:$PATH"' >> "$ZSHRC"
        echo "✓ PATH設定を追加しました"
    else
        echo "✓ PATH設定は既に存在します"
    fi

    # タブ補完設定を追加
    if ! grep -q 'fpath=(~/.zsh/completions $fpath)' "$ZSHRC"; then
        echo "" >> "$ZSHRC"
        echo "# mtg tool - タブ補完設定" >> "$ZSHRC"
        echo 'fpath=(~/.zsh/completions $fpath)' >> "$ZSHRC"
        echo 'autoload -Uz compinit && compinit' >> "$ZSHRC"
        echo "✓ タブ補完設定を追加しました"
    else
        echo "✓ タブ補完設定は既に存在します"
    fi

    echo ""
    echo "✅ 自動セットアップが完了しました！"
    echo ""
    echo "次のステップ:"
    echo "1. シェルを再起動して設定を反映してください:"
    echo "   exec zsh"
    echo ""
    echo "2. 設定ファイルを編集してプロジェクト情報を追加してください:"
    echo "   vim ~/.config/mtg/config.json"
    echo ""
    echo "3. 使い方:"
    echo "   mtg list"
    echo "   mtg prep -project <your-project>"
    echo "   mtg memo -project <your-project>"
else
    echo ""
    echo "⚠️  手動セットアップが必要です"
    echo ""
    echo "次のステップ:"
    echo "1. 設定ファイルを編集してプロジェクト情報を追加してください:"
    echo "   vim ~/.config/mtg/config.json"
    echo ""
    echo "2. ~/.zshrc に以下を追加してください:"
    echo '   export PATH="$HOME/bin:$PATH"'
    echo ""
    echo "3. タブ補完を有効化してください:"
    echo '   fpath=(~/.zsh/completions $fpath)'
    echo '   autoload -Uz compinit && compinit'
    echo ""
    echo "4. 設定を反映してください:"
    echo "   source ~/.zshrc"
    echo ""
    echo "5. 使い方:"
    echo "   mtg list"
    echo "   mtg prep -project <your-project>"
    echo "   mtg memo -project <your-project>"
fi
