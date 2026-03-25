package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/syamaguc/meeting-toolkit/pkg/config"
	"github.com/syamaguc/meeting-toolkit/pkg/file"
	"github.com/syamaguc/meeting-toolkit/pkg/mail"
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
			defer func() { _ = os.Remove(tmpfile.Name()) }()

			if _, err := tmpfile.Write([]byte(tt.configJSON)); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			// テスト実行
			config, err := config.Load(tmpfile.Name())

			if tt.wantErr {
				if err == nil {
					t.Errorf("config.Load() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("config.Load() error = %v, wantErr %v", err, tt.wantErr)
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
	_, err := config.Load("/nonexistent/path/config.json")
	if err == nil {
		t.Error("Load() error = nil, want error for non-existent file")
	}
}

func TestResolvePrefix(t *testing.T) {
	// テスト用の設定ファイルを作成
	tmpfile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()

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
			got, err := config.ResolvePrefix(tt.project, tt.prefix, tt.configPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("config.ResolvePrefix() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("config.ResolvePrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("config.ResolvePrefix() = %v, want %v", got, tt.want)
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
	defer func() { _ = os.RemoveAll(tmpDir) }()

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
	err = file.Rename("PREFIX", tmpDir, currentDate, "")
	if err != nil {
		t.Errorf("Rename() error = %v", err)
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
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// テスト用のファイルを作成
	testFile := "PREFIX_main.txt"
	fpath := filepath.Join(tmpDir, testFile)
	if err := os.WriteFile(fpath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// テスト実行
	currentDate := "20260320"
	suffix := "_MTG後"
	err = file.Rename("PREFIX", tmpDir, currentDate, suffix)
	if err != nil {
		t.Errorf("Rename() error = %v", err)
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
	defer func() { _ = os.RemoveAll(tmpDir) }()

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
	err = file.Collect("PREFIX", tmpDir, destFolder)
	if err != nil {
		t.Errorf("Collect() error = %v", err)
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

func TestLoadConfigWithMailTemplates(t *testing.T) {
	configJSON := `{
		"projects": {
			"test-project": "TEST_PREFIX"
		},
		"mail_templates": {
			"test-project": {
				"prep": "templates/test-prep.txt"
			}
		}
	}`

	tmpfile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	if _, err := tmpfile.Write([]byte(configJSON)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.MailTemplates == nil {
		t.Error("MailTemplates should not be nil")
		return
	}

	projectTemplate, ok := cfg.MailTemplates["test-project"]
	if !ok {
		t.Error("test-project template not found")
		return
	}

	// prep テンプレートパスのチェック
	prepPath, ok := projectTemplate["prep"]
	if !ok {
		t.Error("prep template path not found")
		return
	}
	if prepPath != "templates/test-prep.txt" {
		t.Errorf("prep = %v, want templates/test-prep.txt", prepPath)
	}
}

func TestParseMailTemplate(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantTo      []string
		wantCc      []string
		wantBcc     []string
		wantSubject string
		wantBody    string
		wantErr     bool
	}{
		{
			name: "基本的なテンプレート",
			content: `To: customer@example.com
Cc: team@example.com
Subject: テスト件名

テスト本文
複数行あります`,
			wantTo:      []string{"customer@example.com"},
			wantCc:      []string{"team@example.com"},
			wantBcc:     []string{},
			wantSubject: "テスト件名",
			wantBody:    "テスト本文\n複数行あります",
			wantErr:     false,
		},
		{
			name: "複数のTo/Cc/Bcc",
			content: `To: customer1@example.com, customer2@example.com
Cc: team1@example.com, team2@example.com
Bcc: archive@example.com
Subject: テスト件名

本文`,
			wantTo:      []string{"customer1@example.com", "customer2@example.com"},
			wantCc:      []string{"team1@example.com", "team2@example.com"},
			wantBcc:     []string{"archive@example.com"},
			wantSubject: "テスト件名",
			wantBody:    "本文",
			wantErr:     false,
		},
		{
			name: "件名が空",
			content: `To: customer@example.com
Subject:

本文`,
			wantTo:      []string{"customer@example.com"},
			wantCc:      []string{},
			wantBcc:     []string{},
			wantSubject: "",
			wantBody:    "本文",
			wantErr:     false,
		},
		{
			name: "CcとBccなし",
			content: `To: customer@example.com
Subject: 件名

本文`,
			wantTo:      []string{"customer@example.com"},
			wantCc:      []string{},
			wantBcc:     []string{},
			wantSubject: "件名",
			wantBody:    "本文",
			wantErr:     false,
		},
		{
			name: "スペース含むアドレス",
			content: `To: customer1@example.com , customer2@example.com
Subject: 件名

本文`,
			wantTo:      []string{"customer1@example.com", "customer2@example.com"},
			wantCc:      []string{},
			wantBcc:     []string{},
			wantSubject: "件名",
			wantBody:    "本文",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := mail.Parse(tt.content)

			if tt.wantErr {
				if err == nil {
					t.Error("mail.Parse() error = nil, wantErr true")
				}
				return
			}

			if err != nil {
				t.Errorf("mail.Parse() error = %v, wantErr false", err)
				return
			}

			if !stringSliceEqual(template.To, tt.wantTo) {
				t.Errorf("To = %v, want %v", template.To, tt.wantTo)
			}
			if !stringSliceEqual(template.Cc, tt.wantCc) {
				t.Errorf("Cc = %v, want %v", template.Cc, tt.wantCc)
			}
			if !stringSliceEqual(template.Bcc, tt.wantBcc) {
				t.Errorf("Bcc = %v, want %v", template.Bcc, tt.wantBcc)
			}
			if template.Subject != tt.wantSubject {
				t.Errorf("Subject = %v, want %v", template.Subject, tt.wantSubject)
			}
			if template.Body != tt.wantBody {
				t.Errorf("Body = %v, want %v", template.Body, tt.wantBody)
			}
		})
	}
}

func TestParseMailTemplateWithDate(t *testing.T) {
	currentDate := time.Now().Format("20060102")

	tests := []struct {
		name        string
		content     string
		wantSubject string
		wantBody    string
	}{
		{
			name: "件名に{{DATE}}",
			content: `To: test@example.com
Subject: 【プロジェクト】資料送付 {{DATE}}

本文`,
			wantSubject: "【プロジェクト】資料送付 " + currentDate,
			wantBody:    "本文",
		},
		{
			name: "本文に{{DATE}}",
			content: `To: test@example.com
Subject: 件名

送付資料：
- 資料_{{DATE}}.pdf`,
			wantSubject: "件名",
			wantBody:    "送付資料：\n- 資料_" + currentDate + ".pdf",
		},
		{
			name: "複数の{{DATE}}",
			content: `To: test@example.com
Subject: 資料送付 {{DATE}}

{{DATE}}の資料です。
ファイル名: doc_{{DATE}}.pdf`,
			wantSubject: "資料送付 " + currentDate,
			wantBody:    currentDate + "の資料です。\nファイル名: doc_" + currentDate + ".pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := mail.Parse(tt.content)
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}

			if template.Subject != tt.wantSubject {
				t.Errorf("Subject = %v, want %v", template.Subject, tt.wantSubject)
			}
			if template.Body != tt.wantBody {
				t.Errorf("Body = %v, want %v", template.Body, tt.wantBody)
			}
		})
	}
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestGetMailTemplate(t *testing.T) {
	// テンプレートファイルを作成
	tmpDir, err := os.MkdirTemp("", "mtg-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	templateContent := `To: customer@example.com
Subject: テスト件名

テスト本文`
	templatePath := filepath.Join(tmpDir, "test-prep.txt")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatal(err)
	}

	configJSON := `{
		"projects": {
			"test-project": "TEST_PREFIX"
		},
		"mail_templates": {
			"test-project": {
				"prep": "` + templatePath + `"
			}
		}
	}`

	tmpfile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	if _, err := tmpfile.Write([]byte(configJSON)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		project     string
		mailType    string
		wantErr     bool
		wantSubject string
	}{
		{
			name:        "正常なテンプレート取得",
			project:     "test-project",
			mailType:    "prep",
			wantErr:     false,
			wantSubject: "テスト件名",
		},
		{
			name:     "存在しないプロジェクト",
			project:  "nonexistent",
			mailType: "prep",
			wantErr:  true,
		},
		{
			name:     "不正なメールタイプ",
			project:  "test-project",
			mailType: "invalid",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := mail.Get(tmpfile.Name(), tt.project, tt.mailType)

			if tt.wantErr {
				if err == nil {
					t.Error("getMailTemplate() error = nil, wantErr true")
				}
				return
			}

			if err != nil {
				t.Errorf("getMailTemplate() error = %v, wantErr false", err)
				return
			}

			if template.Subject != tt.wantSubject {
				t.Errorf("Subject = %v, want %v", template.Subject, tt.wantSubject)
			}
		})
	}
}

func TestFormatMailOutput(t *testing.T) {
	tests := []struct {
		name     string
		template *mail.Template
		want     string
	}{
		{
			name: "全フィールド指定",
			template: &mail.Template{
				To:      []string{"customer@example.com", "another@example.com"},
				Cc:      []string{"team@example.com"},
				Bcc:     []string{"bcc@example.com"},
				Subject: "テスト件名",
				Body:    "テスト本文\n複数行あります",
			},
			want: "To: customer@example.com, another@example.com\n" +
				"Cc: team@example.com\n" +
				"Bcc: bcc@example.com\n" +
				"件名: テスト件名\n" +
				"\n" +
				"テスト本文\n" +
				"複数行あります\n",
		},
		{
			name: "To空でCc/Bccなし",
			template: &mail.Template{
				To:      []string{},
				Cc:      []string{},
				Bcc:     []string{},
				Subject: "件名のみ",
				Body:    "本文",
			},
			// To:は必須フィールドとして常に出力、Cc/Bccは空なら出力しない
			want: "To: \n" +
				"件名: 件名のみ\n" +
				"\n" +
				"本文\n",
		},
		{
			name: "件名空",
			template: &mail.Template{
				To:      []string{"to@example.com"},
				Cc:      []string{},
				Bcc:     []string{},
				Subject: "",
				Body:    "本文のみ",
			},
			want: "To: to@example.com\n" +
				"件名: \n" +
				"\n" +
				"本文のみ\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mail.Format(tt.template)
			if got != tt.want {
				t.Errorf("Format() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCreateTemplateFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mtg-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	tests := []struct {
		name         string
		project      string
		mailType     string
		wantFilename string
		wantErr      bool
	}{
		{
			name:         "prep用テンプレート作成",
			project:      "test-project",
			mailType:     "prep",
			wantFilename: "test-project-prep.txt",
			wantErr:      false,
		},
		{
			name:         "memo用テンプレート作成",
			project:      "test-project",
			mailType:     "memo",
			wantFilename: "test-project-memo.txt",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templatePath, existed, err := mail.CreateFile(tmpDir, tt.project, tt.mailType)

			if tt.wantErr {
				if err == nil {
					t.Error("createTemplateFile() error = nil, wantErr true")
				}
				return
			}

			if err != nil {
				t.Errorf("createTemplateFile() error = %v, wantErr false", err)
				return
			}

			if existed {
				t.Error("existed should be false for new template")
			}

			expectedPath := filepath.Join(tmpDir, tt.wantFilename)
			if templatePath != expectedPath {
				t.Errorf("templatePath = %v, want %v", templatePath, expectedPath)
			}

			if _, err := os.Stat(templatePath); os.IsNotExist(err) {
				t.Errorf("Template file not created at %v", templatePath)
			}

			content, err := os.ReadFile(templatePath)
			if err != nil {
				t.Errorf("Failed to read template file: %v", err)
			}

			if len(content) == 0 {
				t.Error("Template file is empty")
			}

			contentStr := string(content)
			if !strings.Contains(contentStr, "To:") {
				t.Error("Template should contain 'To:' header")
			}
			if !strings.Contains(contentStr, "Subject:") {
				t.Error("Template should contain 'Subject:' header")
			}
		})
	}
}

func TestCreateTemplateFileExisting(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mtg-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	filename := "test-project-prep.txt"
	existingPath := filepath.Join(tmpDir, filename)
	existingContent := "existing template content"

	if err := os.WriteFile(existingPath, []byte(existingContent), 0644); err != nil {
		t.Fatal(err)
	}

	templatePath, existed, err := mail.CreateFile(tmpDir, "test-project", "prep")
	if err != nil {
		t.Errorf("CreateFile() should not error when file exists, got: %v", err)
		return
	}

	if !existed {
		t.Error("existed should be true when file already exists")
	}

	if templatePath != existingPath {
		t.Errorf("templatePath = %v, want %v", templatePath, existingPath)
	}

	content, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != existingContent {
		t.Error("Existing file content should be preserved")
	}
}

func TestUpdateConfigWithMailTemplate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mtg-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	tests := []struct {
		name           string
		initialConfig  string
		project        string
		mailType       string
		templatePath   string
		wantErr        bool
		checkFunc      func(*testing.T, *config.Config)
	}{
		{
			name: "新規プロジェクトの追加",
			initialConfig: `{
				"projects": {
					"existing-project": "EXISTING_PREFIX"
				},
				"mail_templates": {}
			}`,
			project:      "new-project",
			mailType:     "prep",
			templatePath: "templates/new-project-prep.txt",
			wantErr:      false,
			checkFunc: func(t *testing.T, cfg *config.Config) {
				if _, ok := cfg.MailTemplates["new-project"]; !ok {
					t.Error("new-project not added to mail_templates")
				}
				if path := cfg.MailTemplates["new-project"]["prep"]; path != "templates/new-project-prep.txt" {
					t.Errorf("prep path = %v, want templates/new-project-prep.txt", path)
				}
			},
		},
		{
			name: "既存プロジェクトに新しいタイプを追加",
			initialConfig: `{
				"projects": {
					"existing-project": "EXISTING_PREFIX"
				},
				"mail_templates": {
					"existing-project": {
						"prep": "templates/existing-prep.txt"
					}
				}
			}`,
			project:      "existing-project",
			mailType:     "memo",
			templatePath: "templates/existing-memo.txt",
			wantErr:      false,
			checkFunc: func(t *testing.T, cfg *config.Config) {
				templates := cfg.MailTemplates["existing-project"]
				if templates["prep"] != "templates/existing-prep.txt" {
					t.Error("Existing prep template should be preserved")
				}
				if templates["memo"] != "templates/existing-memo.txt" {
					t.Errorf("memo path = %v, want templates/existing-memo.txt", templates["memo"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := filepath.Join(tmpDir, "config.json")
			if err := os.WriteFile(configPath, []byte(tt.initialConfig), 0644); err != nil {
				t.Fatal(err)
			}

			err := mail.UpdateConfig(configPath, tt.project, tt.mailType, tt.templatePath)

			if tt.wantErr {
				if err == nil {
					t.Error("updateConfigWithMailTemplate() error = nil, wantErr true")
				}
				return
			}

			if err != nil {
				t.Errorf("updateConfigWithMailTemplate() error = %v, wantErr false", err)
				return
			}

			cfg, err := config.Load(configPath)
			if err != nil {
				t.Fatalf("Failed to load updated config: %v", err)
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, cfg)
			}
		})
	}
}
