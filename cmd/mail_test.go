package cmd

import (
	"os"
	"testing"
)

func TestResolveEditor(t *testing.T) {
	tests := []struct {
		name       string
		envEditor  string
		envVisual  string
		wantEditor string
	}{
		{
			name:       "EDITOR設定あり",
			envEditor:  "vim",
			envVisual:  "",
			wantEditor: "vim",
		},
		{
			name:       "EDITORなし・VISUAL設定あり",
			envEditor:  "",
			envVisual:  "code",
			wantEditor: "code",
		},
		{
			name:       "両方設定なし→viにフォールバック",
			envEditor:  "",
			envVisual:  "",
			wantEditor: "vi",
		},
		{
			name:       "EDITOR優先",
			envEditor:  "nano",
			envVisual:  "code",
			wantEditor: "nano",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 環境変数を保存・復元
			origEditor := os.Getenv("EDITOR")
			origVisual := os.Getenv("VISUAL")
			defer func() {
				_ = os.Setenv("EDITOR", origEditor)
				_ = os.Setenv("VISUAL", origVisual)
			}()

			_ = os.Setenv("EDITOR", tt.envEditor)
			_ = os.Setenv("VISUAL", tt.envVisual)

			got := resolveEditor()
			if got != tt.wantEditor {
				t.Errorf("resolveEditor() = %v, want %v", got, tt.wantEditor)
			}
		})
	}
}
