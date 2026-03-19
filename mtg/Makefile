.PHONY: build install uninstall clean migrate-config

# ビルド設定
BINARY_NAME=mtg
OLD_CONFIG_DIR=$(HOME)/.config/mtg-prep
CONFIG_DIR=$(HOME)/.config/mtg
CONFIG_FILE=$(CONFIG_DIR)/config.json
INSTALL_DIR=$(HOME)/bin

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
	@echo "アンインストール完了!"

clean:
	@rm -f $(BINARY_NAME)
	@echo "ビルド成果物を削除しました"
