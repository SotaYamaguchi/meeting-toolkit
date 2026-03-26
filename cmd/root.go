// Package cmd provides command implementations for the mtg CLI.
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/syamaguc/meeting-toolkit/pkg/config"
)

var rootCmd = &cobra.Command{
	Use:   "mtg",
	Short: "顧客プロジェクトのMTG前後でファイルを整理するツール",
	Long: `mtg - 顧客プロジェクトのMTG前後でファイルを整理するツール

使い方:
  mtg prep [オプション]         MTG前の送付資料を準備
  mtg memo [オプション]         MTG後の議事メモを整理
  mtg mail [オプション]         メールテンプレートを表示
  mtg mail init [オプション]    メールテンプレートを作成
  mtg list                      利用可能なプロジェクト一覧を表示
  mtg completion [shell]        タブ補完スクリプトを出力

例:
  mtg list
  mtg prep -project your-project
  mtg memo -project your-project
  mtg mail -project your-project -type prep
  mtg mail init -project your-project -type prep`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(prepCmd)
	rootCmd.AddCommand(memoCmd)
	rootCmd.AddCommand(mailCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(completionCmd)
}

// completeProjects returns project names from config for flag completion.
func completeProjects(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	cfg, err := config.Load(config.GetDefaultPath())
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var projects []string
	for name := range cfg.Projects {
		projects = append(projects, name)
	}
	return projects, cobra.ShellCompDirectiveNoFileComp
}

// completeMailType returns mail type candidates for flag completion.
func completeMailType(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{"prep", "memo"}, cobra.ShellCompDirectiveNoFileComp
}
