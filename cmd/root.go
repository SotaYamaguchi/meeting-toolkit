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
  mtg files prep [オプション]        MTG前の送付資料を準備
  mtg files post [オプション]        MTG後の議事メモを整理
  mtg mail prep [オプション]         prep用メールテンプレートを表示
  mtg mail post [オプション]         post用メールテンプレートを表示
  mtg mail init prep [オプション]    prep用メールテンプレートを作成
  mtg mail init post [オプション]    post用メールテンプレートを作成
  mtg list                           利用可能なプロジェクト一覧を表示
  mtg completion [shell]             タブ補完スクリプトを出力

例:
  mtg list
  mtg files prep --project your-project
  mtg files post --project your-project
  mtg mail prep --project your-project
  mtg mail init prep --project your-project`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(filesCmd)
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
