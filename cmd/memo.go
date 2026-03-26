package cmd

import (
	"github.com/spf13/cobra"

	"github.com/syamaguc/meeting-toolkit/pkg/config"
	"github.com/syamaguc/meeting-toolkit/pkg/file"
)

var memoCmd = &cobra.Command{
	Use:   "memo",
	Short: "MTG後の議事メモを整理",
	Args: cobra.NoArgs,
	ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		project, _ := cmd.Flags().GetString("project")
		prefix, _ := cmd.Flags().GetString("prefix")
		dir, _ := cmd.Flags().GetString("dir")
		configPath, _ := cmd.Flags().GetString("config")

		finalPrefix, err := config.ResolvePrefix(project, prefix, configPath)
		if err != nil {
			return err
		}

		return file.ProcessMemo(finalPrefix, dir)
	},
}

func init() {
	memoCmd.Flags().StringP("project", "p", "", "プロジェクト名")
	memoCmd.Flags().String("prefix", "", "プレフィックスを直接指定")
	memoCmd.Flags().StringP("dir", "d", ".", "対象ディレクトリ")
	memoCmd.Flags().StringP("config", "c", config.GetDefaultPath(), "設定ファイルのパス")

	_ = memoCmd.RegisterFlagCompletionFunc("project", completeProjects)
}
