// Package main provides a CLI tool for organizing meeting documents.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Projects      map[string]string                     `json:"projects"`
	MailTemplates map[string]map[string]string          `json:"mail_templates"`
}

type MailTemplate struct {
	To      []string `json:"to"`
	Cc      []string `json:"cc"`
	Bcc     []string `json:"bcc"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcommand := os.Args[1]

	switch subcommand {
	case "prep":
		if err := runPrep(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
	case "memo":
		if err := runMemo(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
	case "mail":
		if err := runMail(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
	case "list", "-list", "--list":
		if err := runList(); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		printUsage()
	case "completion":
		if err := runCompletion(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "エラー: 不明なサブコマンド '%s'\n\n", subcommand)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("mtg - 顧客プロジェクトのMTG前後でファイルを整理するツール")
	fmt.Println()
	fmt.Println("使い方:")
	fmt.Println("  mtg prep [オプション]         MTG前の送付資料を準備")
	fmt.Println("  mtg memo [オプション]         MTG後の議事メモを整理")
	fmt.Println("  mtg mail [オプション]         メールテンプレートを表示")
	fmt.Println("  mtg mail init [オプション]    メールテンプレートを作成")
	fmt.Println("  mtg list                      利用可能なプロジェクト一覧を表示")
	fmt.Println("  mtg completion                タブ補完スクリプトを出力")
	fmt.Println()
	fmt.Println("オプション:")
	fmt.Println("  -project <名前>    プロジェクト名 (例: project-a, project-b)")
	fmt.Println("  -prefix <値>       プレフィックスを直接指定")
	fmt.Println("  -dir <パス>        対象ディレクトリ (デフォルト: .)")
	fmt.Println("  -config <パス>     設定ファイルのパス")
	fmt.Println("  -type <値>         メールタイプ (prep または memo, mail サブコマンド用)")
	fmt.Println()
	fmt.Println("例:")
	fmt.Println("  mtg list")
	fmt.Println("  mtg prep -project your-project")
	fmt.Println("  mtg memo -project your-project")
	fmt.Println("  mtg mail -project your-project -type prep")
	fmt.Println("  mtg mail init -project your-project -type prep")
	fmt.Println()
	fmt.Println("タブ補完のセットアップ (zsh):")
	fmt.Println("  mtg completion > ~/.zsh/completions/_mtg")
	fmt.Println("  # ~/.zshrc に以下を追加:")
	fmt.Println("  # fpath=(~/.zsh/completions $fpath)")
	fmt.Println("  # autoload -Uz compinit && compinit")
}

func runList() error {
	defaultConfigPath := getDefaultConfigPath()
	config, err := loadConfig(defaultConfigPath)
	if err != nil {
		return fmt.Errorf("設定ファイル読み込みエラー: %w", err)
	}

	fmt.Println("利用可能なプロジェクト:")
	for proj, pref := range config.Projects {
		fmt.Printf("  %s -> %s\n", proj, pref)
	}
	return nil
}

func runCompletion(args []string) error {
	shell := "zsh"
	if len(args) > 0 {
		shell = args[0]
	}

	switch shell {
	case "zsh":
		fmt.Print(zshCompletionScript)
	case "bash":
		return fmt.Errorf("bash補完は未実装です")
	default:
		return fmt.Errorf("未対応のシェル: %s (zsh のみ対応)", shell)
	}

	return nil
}

const zshCompletionScript = `#compdef mtg

_mtg() {
  local -a subcommands
  subcommands=(
    'prep:MTG前の送付資料を準備'
    'memo:MTG後の議事メモを整理'
    'mail:メールテンプレートを表示'
    'list:利用可能なプロジェクト一覧を表示'
    'completion:タブ補完スクリプトを出力'
    'help:ヘルプを表示'
  )

  local -a options
  options=(
    '-project[プロジェクト名を指定]:project:_mtg_projects'
    '-prefix[プレフィックスを直接指定]:prefix:'
    '-dir[対象ディレクトリを指定]:directory:_files -/'
    '-config[設定ファイルのパスを指定]:config file:_files'
  )

  local -a mail_options
  mail_options=(
    '-project[プロジェクト名を指定]:project:_mtg_projects'
    '-type[メールタイプを指定]:type:(prep memo)'
    '-config[設定ファイルのパスを指定]:config file:_files'
  )

  _arguments -C \
    '1: :->subcommand' \
    '*:: :->args'

  case $state in
    subcommand)
      _describe 'subcommand' subcommands
      ;;
    args)
      case $words[1] in
        prep|memo)
          _arguments $options
          ;;
        mail)
          _arguments $mail_options
          ;;
      esac
      ;;
  esac
}

_mtg_projects() {
  local config_path="$HOME/.config/mtg/config.json"
  if [[ -f "$config_path" ]]; then
    local -a projects
    projects=(${(f)"$(grep -o '"[^"]*":' "$config_path" | tr -d '":' | grep -v -E '(projects|mail_templates|prep|memo)')"})
    _describe 'project' projects
  fi
}

_mtg "$@"
`

func runPrep(args []string) error {
	fs := flag.NewFlagSet("prep", flag.ExitOnError)
	project := fs.String("project", "", "プロジェクト名")
	prefix := fs.String("prefix", "", "プレフィックス")
	dir := fs.String("dir", ".", "対象ディレクトリ")
	configPath := fs.String("config", getDefaultConfigPath(), "設定ファイルのパス")

	if err := fs.Parse(args); err != nil {
		return err
	}

	finalPrefix, err := resolvePrefix(*project, *prefix, *configPath)
	if err != nil {
		return err
	}

	return processPrepFiles(finalPrefix, *dir)
}

func runMemo(args []string) error {
	fs := flag.NewFlagSet("memo", flag.ExitOnError)
	project := fs.String("project", "", "プロジェクト名")
	prefix := fs.String("prefix", "", "プレフィックス")
	dir := fs.String("dir", ".", "対象ディレクトリ")
	configPath := fs.String("config", getDefaultConfigPath(), "設定ファイルのパス")

	if err := fs.Parse(args); err != nil {
		return err
	}

	finalPrefix, err := resolvePrefix(*project, *prefix, *configPath)
	if err != nil {
		return err
	}

	return processMemoFiles(finalPrefix, *dir)
}

func runMail(args []string) error {
	if len(args) > 0 && args[0] == "init" {
		return runMailInit(args[1:])
	}

	fs := flag.NewFlagSet("mail", flag.ExitOnError)
	project := fs.String("project", "", "プロジェクト名")
	mailType := fs.String("type", "", "メールタイプ (prep または memo)")
	configPath := fs.String("config", getDefaultConfigPath(), "設定ファイルのパス")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *project == "" {
		return fmt.Errorf("-project フラグが必要です")
	}

	if *mailType == "" {
		return fmt.Errorf("-type フラグが必要です (prep または memo)")
	}

	template, err := getMailTemplate(*configPath, *project, *mailType)
	if err != nil {
		return err
	}

	output := formatMailOutput(template)
	fmt.Print(output)

	return nil
}

func runMailInit(args []string) error {
	fs := flag.NewFlagSet("mail init", flag.ExitOnError)
	project := fs.String("project", "", "プロジェクト名")
	mailType := fs.String("type", "", "メールタイプ (prep または memo)")
	configPath := fs.String("config", getDefaultConfigPath(), "設定ファイルのパス")

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

	templatePath, existed, err := createTemplateFile(templatesDir, *project, *mailType)
	if err != nil {
		return err
	}

	relPath := "templates/" + filepath.Base(templatePath)

	if err := updateConfigWithMailTemplate(*configPath, *project, *mailType, relPath); err != nil {
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

func resolvePrefix(project, prefix, configPath string) (string, error) {
	if project != "" {
		config, err := loadConfig(configPath)
		if err != nil {
			return "", fmt.Errorf("設定ファイル読み込みエラー: %w", err)
		}

		p, ok := config.Projects[project]
		if !ok {
			fmt.Fprintf(os.Stderr, "エラー: プロジェクト '%s' が見つかりません\n", project)
			fmt.Fprintln(os.Stderr, "\n利用可能なプロジェクト:")
			for proj := range config.Projects {
				fmt.Fprintf(os.Stderr, "  - %s\n", proj)
			}
			return "", fmt.Errorf("プロジェクトが見つかりません")
		}
		return p, nil
	} else if prefix != "" {
		return prefix, nil
	}

	return "", fmt.Errorf("-project または -prefix フラグが必要です")
}

func getDefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "config.json"
	}
	return filepath.Join(home, ".config", "mtg", "config.json")
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ファイル読み込みエラー: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("JSON解析エラー: %w", err)
	}

	return &config, nil
}

func processPrepFiles(prefix, dir string) error {
	currentDate := time.Now().Format("20060102")

	// ファイル名を変更
	if err := renameFiles(prefix, dir, currentDate, ""); err != nil {
		return err
	}

	// フォルダに集約
	destinationFolder := filepath.Join(dir, fmt.Sprintf("%s_資料送付_%s", prefix, currentDate))
	return collectFiles(prefix, dir, destinationFolder)
}

func processMemoFiles(prefix, dir string) error {
	currentDate := time.Now().Format("20060102")
	suffix := "_MTG後"

	// ファイル名を変更
	if err := renameFiles(prefix, dir, currentDate, suffix); err != nil {
		return err
	}

	// フォルダに集約
	destinationFolder := filepath.Join(dir, fmt.Sprintf("%s_資料送付_%s%s", prefix, currentDate, suffix))
	return collectFiles(prefix, dir, destinationFolder)
}

func renameFiles(prefix, dir, currentDate, suffix string) error {
	pattern := filepath.Join(dir, prefix+"*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("ファイル検索エラー: %w", err)
	}

	for _, file := range matches {
		info, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("ファイル情報取得エラー (%s): %w", file, err)
		}
		if info.IsDir() {
			continue
		}

		basename := filepath.Base(file)
		newBasename := strings.ReplaceAll(basename, "main", currentDate+suffix)

		if basename == newBasename {
			continue
		}

		newFile := filepath.Join(dir, newBasename)

		if err := os.Rename(file, newFile); err != nil {
			return fmt.Errorf("ファイル名変更エラー (%s -> %s): %w", file, newFile, err)
		}

		fmt.Printf("- %s\n", newBasename)
	}

	return nil
}

func collectFiles(prefix, dir, destinationFolder string) error {
	if err := os.MkdirAll(destinationFolder, 0755); err != nil {
		return fmt.Errorf("フォルダ作成エラー (%s): %w", destinationFolder, err)
	}

	pattern := filepath.Join(dir, prefix+"_*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("ファイル検索エラー: %w", err)
	}

	for _, file := range matches {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		if info.IsDir() {
			continue
		}

		basename := filepath.Base(file)
		destination := filepath.Join(destinationFolder, basename)

		if err := os.Rename(file, destination); err != nil {
			return fmt.Errorf("ファイル移動エラー (%s -> %s): %w", file, destination, err)
		}
	}

	fmt.Printf("\nファイルを %s に集約しました\n", destinationFolder)
	return nil
}

func getMailTemplate(configPath, project, mailType string) (*MailTemplate, error) {
	config, err := loadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("設定ファイル読み込みエラー: %w", err)
	}

	projectTemplates, ok := config.MailTemplates[project]
	if !ok {
		return nil, fmt.Errorf("プロジェクト '%s' のメールテンプレートが見つかりません", project)
	}

	templatePath, ok := projectTemplates[mailType]
	if !ok {
		return nil, fmt.Errorf("プロジェクト '%s' の %s テンプレートが見つかりません", project, mailType)
	}

	if !filepath.IsAbs(templatePath) {
		configDir := filepath.Dir(configPath)
		templatePath = filepath.Join(configDir, templatePath)
	}

	// テンプレートファイルを読み込み
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("テンプレートファイル読み込みエラー (%s): %w", templatePath, err)
	}

	return parseMailTemplate(string(content))
}

func parseMailTemplate(content string) (*MailTemplate, error) {
	template := &MailTemplate{
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

func formatMailOutput(template *MailTemplate) string {
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

func createTemplateFile(templatesDir, project, mailType string) (string, bool, error) {
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

func updateConfigWithMailTemplate(configPath, project, mailType, templatePath string) error {
	var config *Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config = &Config{
			Projects:      make(map[string]string),
			MailTemplates: make(map[string]map[string]string),
		}
	} else {
		var err error
		config, err = loadConfig(configPath)
		if err != nil {
			return fmt.Errorf("設定ファイル読み込みエラー: %w", err)
		}
	}

	if config.MailTemplates == nil {
		config.MailTemplates = make(map[string]map[string]string)
	}

	if config.MailTemplates[project] == nil {
		config.MailTemplates[project] = make(map[string]string)
	}

	config.MailTemplates[project][mailType] = templatePath

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON変換エラー: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("設定ファイル書き込みエラー: %w", err)
	}

	return nil
}
