package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "タブ補完スクリプトを出力",
	Long: `シェルの補完スクリプトを出力します。

セットアップ (zsh):
  mtg completion zsh > ~/.zsh/completions/_mtg
  # ~/.zshrc に以下を追加:
  # fpath=(~/.zsh/completions $fpath)
  # autoload -Uz compinit && compinit

セットアップ (bash):
  mtg completion bash > /etc/bash_completion.d/mtg`,
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	Args:      cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
	RunE: func(_ *cobra.Command, args []string) error {
		shell := "zsh"
		if len(args) > 0 {
			shell = args[0]
		}

		switch shell {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}

		return nil
	},
}
