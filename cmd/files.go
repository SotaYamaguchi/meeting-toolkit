package cmd

import (
	"github.com/spf13/cobra"

	"github.com/syamaguc/meeting-toolkit/pkg/config"
	"github.com/syamaguc/meeting-toolkit/pkg/file"
)

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "MTG関連ファイルを整理",
}

func newFilesTypeCmd(fileType string) *cobra.Command {
	short := "MTG前の送付資料を準備"
	process := file.ProcessPrep
	if fileType == "post" {
		short = "MTG後の議事メモを整理"
		process = file.ProcessMemo
	}

	cmd := &cobra.Command{
		Use:   fileType,
		Short: short,
		Args:  cobra.NoArgs,
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

			return process(finalPrefix, dir)
		},
	}

	cmd.Flags().StringP("project", "p", "", "プロジェクト名")
	cmd.Flags().String("prefix", "", "プレフィックスを直接指定")
	cmd.Flags().StringP("dir", "d", ".", "対象ディレクトリ")
	cmd.Flags().StringP("config", "c", config.GetDefaultPath(), "設定ファイルのパス")

	_ = cmd.RegisterFlagCompletionFunc("project", completeProjects)

	return cmd
}

func init() {
	filesCmd.AddCommand(newFilesTypeCmd("prep"))
	filesCmd.AddCommand(newFilesTypeCmd("post"))
}
