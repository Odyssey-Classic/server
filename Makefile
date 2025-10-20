SHELL := /bin/bash

# Default ports (can override when invoking make, e.g. `make dev ADMIN_PORT=9001`)
ADMIN_PORT ?= 8081
META_PORT ?= 8082
NETWORK_PORT ?= 8080

# Go parameters
GO_CMD := go
GO_BUILD := $(GO_CMD) build
GO_RUN := $(GO_CMD) run
GO_MOD := $(GO_CMD) mod

# UI parameters
UI_DIR := ui
UI_DEV_CMD := npm run dev --prefix $(UI_DIR)

# Binary output
BIN := odyssey-server

.PHONY: all build server ui dev clean tidy wait-admin open

all: build

build:
	@echo "==> Building server"
	$(GO_BUILD) -o $(BIN) ./cmd

server:
	@echo "==> Starting server (Admin:$(ADMIN_PORT) Meta:$(META_PORT) Network:$(NETWORK_PORT))"
	ODY_ADMIN_PORT=$(ADMIN_PORT) ODY_META_PORT=$(META_PORT) ODY_NETWORK_PORT=$(NETWORK_PORT) $(GO_RUN) ./cmd

ui:
	@echo "==> Starting UI dev server"
	cd $(UI_DIR) && npm install --no-audit --no-fund >/dev/null 2>&1 || true
	$(UI_DEV_CMD)

# Wait until the admin port responds (basic TCP check)
wait-admin:
	@echo "==> Waiting for admin port $(ADMIN_PORT) to be ready..."
	@for i in {1..40}; do \
	  if (echo > /dev/tcp/127.0.0.1/$(ADMIN_PORT)) >/dev/null 2>&1; then \
	    echo "Admin port is up"; \
	    exit 0; \
	  fi; \
	  sleep 0.25; \
	done; \
	echo "Timed out waiting for admin port $(ADMIN_PORT)" >&2; exit 1

# Open the UI in default browser (Linux xdg-open; override if needed)
open:
	@echo "==> Opening UI in browser"
	@xdg-open http://localhost:5173 2>/dev/null || echo "Open http://localhost:5173 manually"

# Run server and UI concurrently. Stops both when interrupted.
dev:
	@echo "==> Starting dev environment"
	@$(MAKE) -j 2 _dev_server _dev_ui

_dev_server:
	@ODY_ADMIN_PORT=$(ADMIN_PORT) ODY_META_PORT=$(META_PORT) ODY_NETWORK_PORT=$(NETWORK_PORT) $(GO_RUN) ./cmd

_dev_ui: wait-admin
	@cd $(UI_DIR) && npm install --no-audit --no-fund >/dev/null 2>&1 || true
	@$(UI_DEV_CMD)

tidy:
	$(GO_MOD) tidy

clean:
	@echo "==> Cleaning build artifacts"
	rm -f $(BIN)

help:
	@echo "Available targets:"
	@echo "  build        Build the Go server binary"
	@echo "  server       Run the server (uses ODY_* env vars or defaults)"
	@echo "  ui           Start UI dev server"
	@echo "  wait-admin   Wait for admin port to accept connections"
	@echo "  dev          Run server + UI concurrently (server first)"
	@echo "  tidy         Run go mod tidy"
	@echo "  clean        Remove built binary"
	@echo "Variables: ADMIN_PORT META_PORT NETWORK_PORT"
