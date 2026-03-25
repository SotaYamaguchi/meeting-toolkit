// Package file provides file operations for meeting document organization.
package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Rename renames files matching the prefix pattern by replacing "main" with the current date.
func Rename(prefix, dir, currentDate, suffix string) error {
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

// Collect moves files matching the prefix pattern to the destination folder.
func Collect(prefix, dir, destinationFolder string) error {
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

// ProcessPrep executes the prep workflow (rename and collect files).
func ProcessPrep(prefix, dir string) error {
	currentDate := time.Now().Format("20060102")

	if err := Rename(prefix, dir, currentDate, ""); err != nil {
		return err
	}

	destinationFolder := filepath.Join(dir, fmt.Sprintf("%s_資料送付_%s", prefix, currentDate))
	return Collect(prefix, dir, destinationFolder)
}

// ProcessMemo executes the memo workflow (rename with suffix and collect files).
func ProcessMemo(prefix, dir string) error {
	currentDate := time.Now().Format("20060102")
	suffix := "_MTG後"

	if err := Rename(prefix, dir, currentDate, suffix); err != nil {
		return err
	}

	destinationFolder := filepath.Join(dir, fmt.Sprintf("%s_資料送付_%s%s", prefix, currentDate, suffix))
	return Collect(prefix, dir, destinationFolder)
}
