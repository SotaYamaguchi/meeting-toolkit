package cmd

import (
	"fmt"

	"github.com/syamaguc/meeting-toolkit/pkg/config"
)

// RunList executes the list subcommand.
func RunList() error {
	cfg, err := config.Load(config.GetDefaultPath())
	if err != nil {
		return fmt.Errorf("設定ファイル読み込みエラー: %w", err)
	}

	fmt.Println("利用可能なプロジェクト:")
	for proj, pref := range cfg.Projects {
		fmt.Printf("  %s -> %s\n", proj, pref)
	}
	return nil
}
