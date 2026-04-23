package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/syamaguc/meeting-toolkit/pkg/config"
	"github.com/syamaguc/meeting-toolkit/pkg/mail"
)

var mailCmd = &cobra.Command{
	Use:   "mail",
	Short: "メールテンプレートを表示・管理",
}

func newMailShowCmd(name, internalType string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("%s用メールテンプレートを表示", name),
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

			template, err := mail.Get(configPath, project, internalType)
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

func newMailInitTypeCmd(name, internalType string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("%s用メールテンプレートを作成", name),
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

			templatePath, existed, err := mail.CreateFile(templatesDir, project, internalType)
			if err != nil {
				return err
			}

			relPath := "templates/" + filepath.Base(templatePath)

			if err := mail.UpdateConfig(configPath, project, internalType, relPath); err != nil {
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

var mailEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "メールテンプレートをエディタで編集",
}

func newMailEditTypeCmd(name, internalType string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("%s用メールテンプレートをエディタで編集", name),
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

			templatePath, err := mail.ResolvePath(configPath, project, internalType)
			if err != nil {
				return err
			}

			if _, err := os.Stat(templatePath); os.IsNotExist(err) {
				return fmt.Errorf("テンプレートファイルが見つかりません: %s\n  'mtg mail init %s -p %s' で作成してください", templatePath, name, project)
			}

			editor := resolveEditor()

			editorCmd := exec.Command(editor, templatePath)
			editorCmd.Stdin = os.Stdin
			editorCmd.Stdout = os.Stdout
			editorCmd.Stderr = os.Stderr

			return editorCmd.Run()
		},
	}

	cmd.Flags().StringP("project", "p", "", "プロジェクト名")
	cmd.Flags().StringP("config", "c", config.GetDefaultPath(), "設定ファイルのパス")

	_ = cmd.RegisterFlagCompletionFunc("project", completeProjects)

	return cmd
}

// resolveEditor returns the editor command to use.
// It checks $EDITOR, $VISUAL, and falls back to "vi".
func resolveEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	return "vi"
}

func init() {
	mailInitCmd.AddCommand(newMailInitTypeCmd("prep", "prep"))
	mailInitCmd.AddCommand(newMailInitTypeCmd("post", "memo"))

	mailEditCmd.AddCommand(newMailEditTypeCmd("prep", "prep"))
	mailEditCmd.AddCommand(newMailEditTypeCmd("post", "memo"))

	mailCmd.AddCommand(newMailShowCmd("prep", "prep"))
	mailCmd.AddCommand(newMailShowCmd("post", "memo"))
	mailCmd.AddCommand(mailInitCmd)
	mailCmd.AddCommand(mailEditCmd)
}
