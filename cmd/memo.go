package cmd

import (
	"flag"

	"github.com/syamaguc/meeting-toolkit/pkg/config"
	"github.com/syamaguc/meeting-toolkit/pkg/file"
)

// RunMemo executes the memo subcommand.
func RunMemo(args []string) error {
	fs := flag.NewFlagSet("memo", flag.ExitOnError)
	project := fs.String("project", "", "プロジェクト名")
	prefix := fs.String("prefix", "", "プレフィックス")
	dir := fs.String("dir", ".", "対象ディレクトリ")
	configPath := fs.String("config", config.GetDefaultPath(), "設定ファイルのパス")

	if err := fs.Parse(args); err != nil {
		return err
	}

	finalPrefix, err := config.ResolvePrefix(*project, *prefix, *configPath)
	if err != nil {
		return err
	}

	return file.ProcessMemo(finalPrefix, *dir)
}
