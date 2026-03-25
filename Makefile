.PHONY: build install uninstall clean lint test help

# ビルド設定
BINARY_NAME=mtg
OLD_CONFIG_DIR=$(HOME)/.config/mtg-prep
CONFIG_DIR=$(HOME)/.config/mtg
CONFIG_FILE=$(CONFIG_DIR)/config.json
INSTALL_DIR=$(HOME)/bin
ZSH_COMPLETION_DIR=$(HOME)/.zsh/completions
ZSH_COMPLETION_FILE=$(ZSH_COMPLETION_DIR)/_mtg

help:
	@echo "mtg - 顧客プロジェクトMTG支援ツール"
	@echo ""
	@echo "使い方:"
	@echo "  make build        ビルド"
	@echo "  make test         テスト実行"
	@echo "  make lint         lint実行"
	@echo "  make install      インストール"
	@echo "  make uninstall    アンインストール"
	@echo "  make clean        クリーンアップ"

build:
	go build -o $(BINARY_NAME)

install: build
	@echo "バイナリをインストール中..."
	@mkdir -p $(INSTALL_DIR)
	@cp $(BINARY_NAME) $(INSTALL_DIR)/
	@echo "設定ファイルをインストール中..."
	@mkdir -p $(CONFIG_DIR)
	@if [ -f $(OLD_CONFIG_DIR)/config.json ] && [ ! -f $(CONFIG_FILE) ]; then \
		echo "既存の設定ファイルを移行中..."; \
		cp $(OLD_CONFIG_DIR)/config.json $(CONFIG_FILE); \
		echo "設定ファイルを $(CONFIG_FILE) に移行しました"; \
	elif [ ! -f $(CONFIG_FILE) ]; then \
		cp config.sample.json $(CONFIG_FILE); \
		echo "設定ファイルを $(CONFIG_FILE) に作成しました"; \
		echo ""; \
		echo "⚠️  重要: $(CONFIG_FILE) を編集してプロジェクト情報を設定してください"; \
	else \
		echo "設定ファイルは既に存在します: $(CONFIG_FILE)"; \
	fi
	@echo ""
	@echo "タブ補完をインストール中..."
	@mkdir -p $(ZSH_COMPLETION_DIR)
	@$(INSTALL_DIR)/$(BINARY_NAME) completion > $(ZSH_COMPLETION_FILE) 2>/dev/null || true
	@if [ -f $(ZSH_COMPLETION_FILE) ]; then \
		echo "zsh補完スクリプトを $(ZSH_COMPLETION_FILE) にインストールしました"; \
		echo ""; \
		echo "⚠️  以下を ~/.zshrc に追加して、シェルを再起動してください:"; \
		echo "  fpath=(~/.zsh/completions \$$fpath)"; \
		echo "  autoload -Uz compinit && compinit"; \
	fi
	@echo ""
	@echo "古いツールを削除しています..."
	@rm -f $(INSTALL_DIR)/mtg-prep $(INSTALL_DIR)/mtg-memo
	@echo ""
	@echo "インストール完了!"
	@echo "$(INSTALL_DIR) がPATHに含まれていることを確認してください"
	@echo ""
	@echo "使い方:"
	@echo "  $(BINARY_NAME) list"
	@echo "  $(BINARY_NAME) prep -project <your-project>"
	@echo "  $(BINARY_NAME) memo -project <your-project>"

uninstall:
	@echo "バイナリを削除中..."
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "設定ファイルを削除中..."
	@rm -rf $(CONFIG_DIR)
	@echo "補完スクリプトを削除中..."
	@rm -f $(ZSH_COMPLETION_FILE)
	@echo "アンインストール完了!"

clean:
	@rm -f $(BINARY_NAME) coverage.out
	@echo "ビルド成果物を削除しました"

lint:
	@echo "Linting..."
	@golangci-lint run

test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
