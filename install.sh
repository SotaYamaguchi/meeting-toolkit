#!/bin/bash
set -e

echo "mtg - 顧客プロジェクトMTG支援ツールをインストールします"
echo ""

# mtgディレクトリに移動してインストール
cd "$(dirname "$0")/mtg"
make install

echo ""
echo "インストールが完了しました！"
echo ""
echo "次のステップ:"
echo "1. 設定ファイルを編集してプロジェクト情報を追加してください:"
echo "   vim ~/.config/mtg/config.json"
echo ""
echo "2. ~/.zshrc または ~/.bashrc に以下を追加してください:"
echo "   export PATH=\"\$HOME/bin:\$PATH\""
echo ""
echo "3. 設定を反映してください:"
echo "   source ~/.zshrc  # または source ~/.bashrc"
echo ""
echo "4. 使い方:"
echo "   mtg list"
echo "   mtg prep -project <your-project>"
echo "   mtg memo -project <your-project>"
