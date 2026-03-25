package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/syamaguc/meeting-toolkit/pkg/config"
	"github.com/syamaguc/meeting-toolkit/pkg/mail"
)

// RunMail executes the mail subcommand.
func RunMail(args []string) error {
	if len(args) > 0 && args[0] == "init" {
		return RunMailInit(args[1:])
	}

	fs := flag.NewFlagSet("mail", flag.ExitOnError)
	project := fs.String("project", "", "プロジェクト名")
	mailType := fs.String("type", "", "メールタイプ (prep または memo)")
	configPath := fs.String("config", config.GetDefaultPath(), "設定ファイルのパス")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *project == "" {
		return fmt.Errorf("-project フラグが必要です")
	}

	if *mailType == "" {
		return fmt.Errorf("-type フラグが必要です (prep または memo)")
	}

	template, err := mail.Get(*configPath, *project, *mailType)
	if err != nil {
		return err
	}

	output := mail.Format(template)
	fmt.Print(output)

	return nil
}

// RunMailInit executes the mail init subcommand.
func RunMailInit(args []string) error {
	fs := flag.NewFlagSet("mail init", flag.ExitOnError)
	project := fs.String("project", "", "プロジェクト名")
	mailType := fs.String("type", "", "メールタイプ (prep または memo)")
	configPath := fs.String("config", config.GetDefaultPath(), "設定ファイルのパス")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *project == "" {
		return fmt.Errorf("-project フラグが必要です")
	}

	if *mailType == "" {
		return fmt.Errorf("-type フラグが必要です (prep または memo)")
	}

	if *mailType != "prep" && *mailType != "memo" {
		return fmt.Errorf("不正なメールタイプ: %s (prep または memo を指定してください)", *mailType)
	}

	configDir := filepath.Dir(*configPath)
	templatesDir := filepath.Join(configDir, "templates")

	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return fmt.Errorf("テンプレートディレクトリ作成エラー: %w", err)
	}

	templatePath, existed, err := mail.CreateFile(templatesDir, *project, *mailType)
	if err != nil {
		return err
	}

	relPath := "templates/" + filepath.Base(templatePath)

	if err := mail.UpdateConfig(*configPath, *project, *mailType, relPath); err != nil {
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
}
