# このリポジトリについて
stephenafamo/bobを使ってcode generateするリポジトリです。　　
- 動作環境としては以下のツールのインストールが必要になります
  - docker
  - Makefile
  - direnv

インストール後に以下の通りに使用してください.
使用方法としては、以下の通りです。
- `pkg/gen/models`のdirを削除。
- `make up`でdocker compose upを行い、mysqlサーバーの立ち上げとマイグレーションを行います。
- マイグレーション後に`go run github.com/stephenafamo/bob/gen/bobgen-mysql@latest -c ./bobgen.yaml`を実行し、`pkg/gen/models`にコードを生成します。
- 
