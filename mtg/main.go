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
