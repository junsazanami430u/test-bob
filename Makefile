SHELL=/bin/bash

.DEFAULT_GOAL := help

.PHONY: help

help: ## show this message
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-8s\033[0m %s\n", $$1, $$2}'

local: ## ko build (ko.local)
	$(MAKE) publish KO_DOCKER_REPO=ko.local

fmt: ## go fmt
	fmt ./...

lint: ## golangci-lint
	golangci-lint run -c golangci.yml --timeout 10m

test: ## go test
test: up
	# パラレルにテストを実行する
	go test -v ./...

govulncheck: ## go vuln check
	govulncheck ./...

run: ## go run .
	-@go run .

deps: ## install dependencies
	env GOBIN=$(PWD)/tmp/bin go install github.com/volatiletech/sqlboiler/v4@v4.18.0
	env GOBIN=$(PWD)/tmp/bin go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@v4.18.0
	env GOBIN=$(PWD)/tmp/bin go install github.com/rubenv/sql-migrate/sql-migrate@v1.7.1
	env GOBIN=$(PWD)/tmp/bin go install github.com/k1LoW/tbls@v1.81.0
	env GOBIN=$(PWD)/tmp/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5
	env GOBIN=$(PWD)/tmp/bin go install golang.org/x/vuln/cmd/govulncheck@latest
	env GOBIN=$(PWD)/tmp/bin go install go.uber.org/mock/mockgen@v0.5.0

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
