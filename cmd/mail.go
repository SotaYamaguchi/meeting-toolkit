package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/syamaguc/meeting-toolkit/pkg/config"
	"github.com/syamaguc/meeting-toolkit/pkg/mail"
)

var mailCmd = &cobra.Command{
	Use:   "mail",
	Short: "メールテンプレートを表示・管理",
}

func newMailShowCmd(mailType string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   mailType,
		Short: fmt.Sprintf("%s用メールテンプレートを表示", mailType),
		Args:  cobra.NoArgs,
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			project, _ := cmd.Flags().GetString("project")
			configPath, _ := cmd.Flags().GetString("config")

			if project == "" {
				return fmt.Errorf("--project フラグが必要です")
			}

			template, err := mail.Get(configPath, project, mailType)
			if err != nil {
				return err
			}

			output := mail.Format(template)
			fmt.Print(output)

			return nil
		},
	}

	cmd.Flags().StringP("project", "p", "", "プロジェクト名")
	cmd.Flags().StringP("config", "c", config.GetDefaultPath(), "設定ファイルのパス")

	_ = cmd.RegisterFlagCompletionFunc("project", completeProjects)

	return cmd
}

var mailInitCmd = &cobra.Command{
	Use:   "init",
	Short: "メールテンプレートファイルを作成",
}

func newMailInitTypeCmd(mailType string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   mailType,
		Short: fmt.Sprintf("%s用メールテンプレートを作成", mailType),
		Args:  cobra.NoArgs,
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			project, _ := cmd.Flags().GetString("project")
			configPath, _ := cmd.Flags().GetString("config")

			if project == "" {
				return fmt.Errorf("--project フラグが必要です")
			}

			configDir := filepath.Dir(configPath)
			templatesDir := filepath.Join(configDir, "templates")

			if err := os.MkdirAll(templatesDir, 0755); err != nil {
				return fmt.Errorf("テンプレートディレクトリ作成エラー: %w", err)
			}

			templatePath, existed, err := mail.CreateFile(templatesDir, project, mailType)
			if err != nil {
				return err
			}

			relPath := "templates/" + filepath.Base(templatePath)

			if err := mail.UpdateConfig(configPath, project, mailType, relPath); err != nil {
				return err
			}

			if existed {
				fmt.Printf("⚠️  テンプレートファイルは既に存在します: %s\n", templatePath)
				fmt.Printf("✓ 既存のファイルを使用します\n")
			} else {
				fmt.Printf("✓ テンプレートファイルを作成しました: %s\n", templatePath)
			}
			fmt.Printf("✓ config.jsonを更新しました\n")
			fmt.Printf("\nテンプレートを編集してください:\n")
			fmt.Printf("  vim %s\n", templatePath)

			return nil
		},
	}

	cmd.Flags().StringP("project", "p", "", "プロジェクト名")
	cmd.Flags().StringP("config", "c", config.GetDefaultPath(), "設定ファイルのパス")

	_ = cmd.RegisterFlagCompletionFunc("project", completeProjects)

	return cmd
}

func init() {
	mailInitCmd.AddCommand(newMailInitTypeCmd("prep"))
	mailInitCmd.AddCommand(newMailInitTypeCmd("memo"))

	mailCmd.AddCommand(newMailShowCmd("prep"))
	mailCmd.AddCommand(newMailShowCmd("memo"))
	mailCmd.AddCommand(mailInitCmd)
}
