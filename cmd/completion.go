package cmd

import (
	"fmt"
)

// RunCompletion executes the completion subcommand.
func RunCompletion(args []string) error {
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

  local -a mail_subcommands
  mail_subcommands=(
    'init:メールテンプレートファイルを作成'
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
          _arguments -C \
            '1: :->mail_subcommand' \
            '*:: :->mail_args'
          case $state in
            mail_subcommand)
              _describe 'mail subcommand' mail_subcommands
              _arguments $mail_options
              ;;
            mail_args)
              case $words[1] in
                init)
                  _arguments $mail_options
                  ;;
              esac
              ;;
          esac
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
