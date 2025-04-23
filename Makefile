SHELL=/bin/bash

.DEFAULT_GOAL := help

.PHONY: help

help: ## show this message
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-8s\033[0m %s\n", $$1, $$2}'

fmt: ## go fmt
	fmt ./...

test: ## go test
test: up
	# パラレルにテストを実行する
	go test -v ./...

govulncheck: ## go vuln check
	govulncheck ./...

run: ## go run .
	-@go run .
up:	## up local 開発に必要なコンテナを上げる
	docker compose up -d

down:	## down local開発に必要なコンテナを落として。MySQLのでデータも削除する。
	docker compose down
	rm -rf volumes
	rm -rf database/atlas.sum

## 複数モジュールに渡って作業するので、必要なところに、go workを作成する
gowork:	## go.work を作る
gowork: test-bob-gowork # cli-gowork core-migrate-gowork  pkg-gowork mocpkg-gowork
	@echo "go work を作成しました"

## patient-assignment のgo.workを作成
test-bob-gowork:
	@- [ ! -f go.work ] && go work init . || true

gowork-clean:	## go.work を削除する
	rm go.work

env: ## 環境変数を読み込む
	direnv allow
