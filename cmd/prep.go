package cmd

import (
	"github.com/spf13/cobra"

	"github.com/syamaguc/meeting-toolkit/pkg/config"
	"github.com/syamaguc/meeting-toolkit/pkg/file"
)

var prepCmd = &cobra.Command{
	Use:   "prep",
	Short: "MTG前の送付資料を準備",
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

		return file.ProcessPrep(finalPrefix, dir)
	},
}

func init() {
	prepCmd.Flags().StringP("project", "p", "", "プロジェクト名")
	prepCmd.Flags().String("prefix", "", "プレフィックスを直接指定")
	prepCmd.Flags().StringP("dir", "d", ".", "対象ディレクトリ")
	prepCmd.Flags().StringP("config", "c", config.GetDefaultPath(), "設定ファイルのパス")

	_ = prepCmd.RegisterFlagCompletionFunc("project", completeProjects)
}
