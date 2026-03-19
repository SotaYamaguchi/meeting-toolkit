.PHONY: install uninstall build clean help

help:
	@echo "mtg - 顧客プロジェクトMTG支援ツール"
	@echo ""
	@echo "使い方:"
	@echo "  make install      インストール"
	@echo "  make uninstall    アンインストール"
	@echo "  make build        ビルドのみ"
	@echo "  make clean        クリーンアップ"

install:
	@cd mtg && $(MAKE) install

uninstall:
	@cd mtg && $(MAKE) uninstall

build:
	@cd mtg && $(MAKE) build

clean:
	@cd mtg && $(MAKE) clean
