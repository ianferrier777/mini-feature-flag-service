#!/bin/bash

set -e

PROJECT_NAME="mini-feature-flags"
MODULE_PATH="github.com/yourusername/$PROJECT_NAME"

echo "📦 Initializing root module..."
go mod init $MODULE_PATH

echo "📂 Creating dummy files for tidy resolution..."
# Make sure each Go file exists (if not already created manually)
touch cmd/server/main.go
touch internal/api/handler.go
touch internal/config/config.go
touch internal/flags/service.go
touch internal/flags/model.go
touch internal/flags/store.go

echo "🔧 Adding replace paths to go.mod..."
# Append local replace directives manually (Go 1.18+ supports this pattern)
cat <<EOL >> go.mod

replace $MODULE_PATH/internal/api => ./internal/api
replace $MODULE_PATH/internal/config => ./internal/config
replace $MODULE_PATH/internal/flags => ./internal/flags
EOL

echo "🧹 Running go mod tidy..."
go mod tidy

echo "✅ Go modules initialized and tidy complete."