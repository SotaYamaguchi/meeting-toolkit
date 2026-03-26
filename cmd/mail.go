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
	Short: "メールテンプレートを表示",
	RunE: func(cmd *cobra.Command, _ []string) error {
		project, _ := cmd.Flags().GetString("project")
		mailType, _ := cmd.Flags().GetString("type")
		configPath, _ := cmd.Flags().GetString("config")

		if project == "" {
			return fmt.Errorf("-project フラグが必要です")
		}

		if mailType == "" {
			return fmt.Errorf("-type フラグが必要です (prep または memo)")
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

var mailInitCmd = &cobra.Command{
	Use:   "init",
	Short: "メールテンプレートファイルを作成",
	RunE: func(cmd *cobra.Command, _ []string) error {
		project, _ := cmd.Flags().GetString("project")
		mailType, _ := cmd.Flags().GetString("type")
		configPath, _ := cmd.Flags().GetString("config")

		if project == "" {
			return fmt.Errorf("-project フラグが必要です")
		}

		if mailType == "" {
			return fmt.Errorf("-type フラグが必要です (prep または memo)")
		}

		if mailType != "prep" && mailType != "memo" {
			return fmt.Errorf("不正なメールタイプ: %s (prep または memo を指定してください)", mailType)
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

func init() {
	mailCmd.Flags().StringP("project", "p", "", "プロジェクト名")
	mailCmd.Flags().StringP("type", "t", "", "メールタイプ (prep または memo)")
	mailCmd.Flags().StringP("config", "c", config.GetDefaultPath(), "設定ファイルのパス")

	mailInitCmd.Flags().StringP("project", "p", "", "プロジェクト名")
	mailInitCmd.Flags().StringP("type", "t", "", "メールタイプ (prep または memo)")
	mailInitCmd.Flags().StringP("config", "c", config.GetDefaultPath(), "設定ファイルのパス")

	mailCmd.AddCommand(mailInitCmd)
}
