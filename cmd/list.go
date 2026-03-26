package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/syamaguc/meeting-toolkit/pkg/config"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "利用可能なプロジェクト一覧を表示",
	RunE: func(_ *cobra.Command, _ []string) error {
		cfg, err := config.Load(config.GetDefaultPath())
		if err != nil {
			return fmt.Errorf("設定ファイル読み込みエラー: %w", err)
		}

		fmt.Println("利用可能なプロジェクト:")
		for proj, pref := range cfg.Projects {
			fmt.Printf("  %s -> %s\n", proj, pref)
		}
		return nil
	},
}
