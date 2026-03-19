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
	Projects map[string]string `json:"projects"`
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
	fmt.Println("  mtg prep [オプション]    MTG前の送付資料を準備")
	fmt.Println("  mtg memo [オプション]    MTG後の議事メモを整理")
	fmt.Println("  mtg list                 利用可能なプロジェクト一覧を表示")
	fmt.Println("  mtg completion           タブ補完スクリプトを出力")
	fmt.Println()
	fmt.Println("オプション:")
	fmt.Println("  -project <名前>    プロジェクト名 (例: project-a, project-b)")
	fmt.Println("  -prefix <値>       プレフィックスを直接指定")
	fmt.Println("  -dir <パス>        対象ディレクトリ (デフォルト: .)")
	fmt.Println("  -config <パス>     設定ファイルのパス")
	fmt.Println()
	fmt.Println("例:")
	fmt.Println("  mtg list")
	fmt.Println("  mtg prep -project your-project")
	fmt.Println("  mtg memo -project your-project")
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
      esac
      ;;
  esac
}

_mtg_projects() {
  local config_path="$HOME/.config/mtg/config.json"
  if [[ -f "$config_path" ]]; then
    local -a projects
    projects=(${(f)"$(grep -o '"[^"]*":' "$config_path" | tr -d '":' | grep -v projects)"})
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
