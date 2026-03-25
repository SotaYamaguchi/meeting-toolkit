// Package cmd provides command implementations for the mtg CLI.
package cmd

import (
	"fmt"
)

// PrintUsage prints the usage information for the mtg command.
func PrintUsage() {
	fmt.Println("mtg - 顧客プロジェクトのMTG前後でファイルを整理するツール")
	fmt.Println()
	fmt.Println("使い方:")
	fmt.Println("  mtg prep [オプション]         MTG前の送付資料を準備")
	fmt.Println("  mtg memo [オプション]         MTG後の議事メモを整理")
	fmt.Println("  mtg mail [オプション]         メールテンプレートを表示")
	fmt.Println("  mtg mail init [オプション]    メールテンプレートを作成")
	fmt.Println("  mtg list                      利用可能なプロジェクト一覧を表示")
	fmt.Println("  mtg completion                タブ補完スクリプトを出力")
	fmt.Println()
	fmt.Println("オプション:")
	fmt.Println("  -project <名前>    プロジェクト名 (例: project-a, project-b)")
	fmt.Println("  -prefix <値>       プレフィックスを直接指定")
	fmt.Println("  -dir <パス>        対象ディレクトリ (デフォルト: .)")
	fmt.Println("  -config <パス>     設定ファイルのパス")
	fmt.Println("  -type <値>         メールタイプ (prep または memo, mail サブコマンド用)")
	fmt.Println()
	fmt.Println("例:")
	fmt.Println("  mtg list")
	fmt.Println("  mtg prep -project your-project")
	fmt.Println("  mtg memo -project your-project")
	fmt.Println("  mtg mail -project your-project -type prep")
	fmt.Println("  mtg mail init -project your-project -type prep")
	fmt.Println()
	fmt.Println("タブ補完のセットアップ (zsh):")
	fmt.Println("  mtg completion > ~/.zsh/completions/_mtg")
	fmt.Println("  # ~/.zshrc に以下を追加:")
	fmt.Println("  # fpath=(~/.zsh/completions $fpath)")
	fmt.Println("  # autoload -Uz compinit && compinit")
}
