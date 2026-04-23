// Package mail provides email template management.
package mail

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/syamaguc/meeting-toolkit/pkg/config"
)

// Template represents an email template.
type Template struct {
	// To contains recipient email addresses.
	To []string
	// Cc contains carbon copy email addresses.
	Cc []string
	// Bcc contains blind carbon copy email addresses.
	Bcc []string
	// Subject is the email subject line.
	Subject string
	// Body is the email body content.
	Body string
}

// ResolvePath resolves the absolute path to a mail template file for a project and mail type.
func ResolvePath(configPath, project, mailType string) (string, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return "", fmt.Errorf("設定ファイル読み込みエラー: %w", err)
	}

	projectTemplates, ok := cfg.MailTemplates[project]
	if !ok {
		return "", fmt.Errorf("プロジェクト '%s' のメールテンプレートが見つかりません", project)
	}

	templatePath, ok := projectTemplates[mailType]
	if !ok {
		return "", fmt.Errorf("プロジェクト '%s' の %s テンプレートが見つかりません", project, mailType)
	}

	if !filepath.IsAbs(templatePath) {
		configDir := filepath.Dir(configPath)
		templatePath = filepath.Join(configDir, templatePath)
	}

	return templatePath, nil
}

// Get retrieves a mail template for a project and mail type.
func Get(configPath, project, mailType string) (*Template, error) {
	templatePath, err := ResolvePath(configPath, project, mailType)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("テンプレートファイル読み込みエラー (%s): %w", templatePath, err)
	}

	return Parse(string(content))
}

// Parse parses template content into a Template struct.
func Parse(content string) (*Template, error) {
	template := &Template{
		To:  []string{},
		Cc:  []string{},
		Bcc: []string{},
	}

	currentDate := time.Now().Format("20060102")

	lines := strings.Split(content, "\n")
	bodyStart := -1

	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			bodyStart = i + 1
			break
		}

		if addresses, found := strings.CutPrefix(line, "To:"); found {
			template.To = parseEmailAddresses(addresses)
		} else if addresses, found := strings.CutPrefix(line, "Cc:"); found {
			template.Cc = parseEmailAddresses(addresses)
		} else if addresses, found := strings.CutPrefix(line, "Bcc:"); found {
			template.Bcc = parseEmailAddresses(addresses)
		} else if subject, found := strings.CutPrefix(line, "Subject:"); found {
			template.Subject = strings.TrimSpace(subject)
		}
	}

	if bodyStart >= 0 && bodyStart < len(lines) {
		template.Body = strings.Join(lines[bodyStart:], "\n")
	}

	template.Subject = strings.ReplaceAll(template.Subject, "{{DATE}}", currentDate)
	template.Body = strings.ReplaceAll(template.Body, "{{DATE}}", currentDate)

	return template, nil
}

// Format formats a template for output.
func Format(template *Template) string {
	var output strings.Builder

	output.WriteString("To: ")
	if len(template.To) > 0 {
		output.WriteString(strings.Join(template.To, ", "))
	}
	output.WriteString("\n")

	if len(template.Cc) > 0 {
		output.WriteString("Cc: ")
		output.WriteString(strings.Join(template.Cc, ", "))
		output.WriteString("\n")
	}

	if len(template.Bcc) > 0 {
		output.WriteString("Bcc: ")
		output.WriteString(strings.Join(template.Bcc, ", "))
		output.WriteString("\n")
	}

	output.WriteString("件名: ")
	output.WriteString(template.Subject)
	output.WriteString("\n")

	output.WriteString("\n")
	output.WriteString(template.Body)
	output.WriteString("\n")

	return output.String()
}

// CreateFile creates a template file with default content.
func CreateFile(templatesDir, project, mailType string) (string, bool, error) {
	filename := fmt.Sprintf("%s-%s.txt", project, mailType)
	templatePath := filepath.Join(templatesDir, filename)

	if _, err := os.Stat(templatePath); err == nil {
		return templatePath, true, nil
	}

	defaultTemplate := `To:
Cc:
Subject:

メール本文をここに記入してください。

`

	if err := os.WriteFile(templatePath, []byte(defaultTemplate), 0644); err != nil {
		return "", false, fmt.Errorf("テンプレートファイル作成エラー: %w", err)
	}

	return templatePath, false, nil
}

// UpdateConfig updates the configuration file with a mail template path.
func UpdateConfig(configPath, project, mailType, templatePath string) error {
	var cfg *config.Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg = &config.Config{
			Projects:      make(map[string]string),
			MailTemplates: make(map[string]map[string]string),
		}
	} else {
		var err error
		cfg, err = config.Load(configPath)
		if err != nil {
			return fmt.Errorf("設定ファイル読み込みエラー: %w", err)
		}
	}

	if cfg.MailTemplates == nil {
		cfg.MailTemplates = make(map[string]map[string]string)
	}

	if cfg.MailTemplates[project] == nil {
		cfg.MailTemplates[project] = make(map[string]string)
	}

	cfg.MailTemplates[project][mailType] = templatePath

	return config.Save(configPath, cfg)
}

func parseEmailAddresses(addressLine string) []string {
	if strings.TrimSpace(addressLine) == "" {
		return []string{}
	}

	addresses := strings.Split(addressLine, ",")
	result := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		trimmed := strings.TrimSpace(addr)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
