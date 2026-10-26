BIN_DIR := build
CMDS := server console scheduler worker

.PHONY: build clean dev migrate

build:
	@for c in $(CMDS); do go build -o $(BIN_DIR)/$$c ./cmd/$$c || exit 1; done

clean:
	find $(BIN_DIR) -type f ! -name .gitkeep -delete

# マイグレーション（make migrate ARGS=down / ARGS=status）
migrate:
	go run ./cmd/console migrate $(ARGS)

# ホットリロード（air未導入なら go run で代用）
dev:
	@if command -v air >/dev/null 2>&1; then air; else go run github.com/air-verse/air@latest; fi
