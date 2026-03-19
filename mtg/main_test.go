package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		configJSON  string
		wantErr     bool
		wantProject map[string]string
	}{
		{
			name: "正常なconfig",
			configJSON: `{
				"projects": {
					"project-a": "PREFIX_A",
					"project-b": "PREFIX_B"
				}
			}`,
			wantErr: false,
			wantProject: map[string]string{
				"project-a": "PREFIX_A",
				"project-b": "PREFIX_B",
			},
		},
		{
			name:       "不正なJSON",
			configJSON: `{"projects": {`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 一時ファイルを作成
			tmpfile, err := os.CreateTemp("", "config*.json")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.configJSON)); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			// テスト実行
			config, err := loadConfig(tmpfile.Name())

			if tt.wantErr {
				if err == nil {
					t.Errorf("loadConfig() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(config.Projects) != len(tt.wantProject) {
				t.Errorf("projects count = %v, want %v", len(config.Projects), len(tt.wantProject))
				return
			}

			for key, want := range tt.wantProject {
				if got, ok := config.Projects[key]; !ok || got != want {
					t.Errorf("projects[%s] = %v, want %v", key, got, want)
				}
			}
		})
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := loadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Error("loadConfig() error = nil, want error for non-existent file")
	}
}

func TestResolvePrefix(t *testing.T) {
	// テスト用の設定ファイルを作成
	tmpfile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	configData := map[string]any{
		"projects": map[string]string{
			"test-project": "TEST_PREFIX",
		},
	}
	data, err := json.Marshal(configData)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		project    string
		prefix     string
		configPath string
		want       string
		wantErr    bool
	}{
		{
			name:       "プロジェクト名で解決",
			project:    "test-project",
			prefix:     "",
			configPath: tmpfile.Name(),
			want:       "TEST_PREFIX",
			wantErr:    false,
		},
		{
			name:       "プレフィックスを直接指定",
			project:    "",
			prefix:     "DIRECT_PREFIX",
			configPath: tmpfile.Name(),
			want:       "DIRECT_PREFIX",
			wantErr:    false,
		},
		{
			name:       "存在しないプロジェクト",
			project:    "nonexistent",
			prefix:     "",
			configPath: tmpfile.Name(),
			want:       "",
			wantErr:    true,
		},
		{
			name:       "project と prefix 両方未指定",
			project:    "",
			prefix:     "",
			configPath: tmpfile.Name(),
			want:       "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolvePrefix(tt.project, tt.prefix, tt.configPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("resolvePrefix() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("resolvePrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("resolvePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenameFiles(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir, err := os.MkdirTemp("", "mtg-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// テスト用のファイルを作成
	testFiles := []string{
		"PREFIX_main.txt",
		"PREFIX_main_document.pdf",
		"PREFIX_other.txt",
	}

	for _, fname := range testFiles {
		fpath := filepath.Join(tmpDir, fname)
		if err := os.WriteFile(fpath, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// テスト実行
	currentDate := "20260320"
	err = renameFiles("PREFIX", tmpDir, currentDate, "")
	if err != nil {
		t.Errorf("renameFiles() error = %v", err)
		return
	}

	// 結果確認
	expectedFiles := map[string]bool{
		"PREFIX_20260320.txt":          true,
		"PREFIX_20260320_document.pdf": true,
		"PREFIX_other.txt":             true,
	}

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	foundFiles := make(map[string]bool)
	for _, entry := range entries {
		foundFiles[entry.Name()] = true
	}

	for fname := range expectedFiles {
		if !foundFiles[fname] {
			t.Errorf("Expected file %s not found", fname)
		}
	}

	for fname := range foundFiles {
		if !expectedFiles[fname] {
			t.Errorf("Unexpected file %s found", fname)
		}
	}
}

func TestRenameFilesWithSuffix(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir, err := os.MkdirTemp("", "mtg-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// テスト用のファイルを作成
	testFile := "PREFIX_main.txt"
	fpath := filepath.Join(tmpDir, testFile)
	if err := os.WriteFile(fpath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// テスト実行
	currentDate := "20260320"
	suffix := "_MTG後"
	err = renameFiles("PREFIX", tmpDir, currentDate, suffix)
	if err != nil {
		t.Errorf("renameFiles() error = %v", err)
		return
	}

	// 結果確認
	expectedFile := "PREFIX_20260320_MTG後.txt"
	if _, err := os.Stat(filepath.Join(tmpDir, expectedFile)); os.IsNotExist(err) {
		t.Errorf("Expected file %s not found", expectedFile)
	}
}

func TestCollectFiles(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir, err := os.MkdirTemp("", "mtg-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// テスト用のファイルを作成
	testFiles := []string{
		"PREFIX_20260320.txt",
		"PREFIX_20260320_document.pdf",
		"OTHER_file.txt",
	}

	for _, fname := range testFiles {
		fpath := filepath.Join(tmpDir, fname)
		if err := os.WriteFile(fpath, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// テスト実行
	destFolder := filepath.Join(tmpDir, "PREFIX_資料送付_20260320")
	err = collectFiles("PREFIX", tmpDir, destFolder)
	if err != nil {
		t.Errorf("collectFiles() error = %v", err)
		return
	}

	// 結果確認: 移動先ディレクトリが作成されたか
	if _, err := os.Stat(destFolder); os.IsNotExist(err) {
		t.Errorf("Destination folder %s not created", destFolder)
		return
	}

	// 結果確認: ファイルが移動されたか
	expectedInDest := []string{
		"PREFIX_20260320.txt",
		"PREFIX_20260320_document.pdf",
	}

	for _, fname := range expectedInDest {
		destPath := filepath.Join(destFolder, fname)
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			t.Errorf("Expected file %s not found in destination", fname)
		}

		// 元の場所には存在しないことを確認
		srcPath := filepath.Join(tmpDir, fname)
		if _, err := os.Stat(srcPath); !os.IsNotExist(err) {
			t.Errorf("File %s should have been moved from source", fname)
		}
	}

	// OTHER_file.txt は移動されていないことを確認
	otherPath := filepath.Join(tmpDir, "OTHER_file.txt")
	if _, err := os.Stat(otherPath); os.IsNotExist(err) {
		t.Errorf("File OTHER_file.txt should not have been moved")
	}
}
