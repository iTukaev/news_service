.PHONY: receiver validator data mailing client
service: s_build
	@./service
s_build:
	@go build -o service ./cmd/service/*.go

client: c_build
	@./client
c_build:
	@go build -o client ./cmd/client/*.go


MIGRATION_DIR:=./migrations
.PHONY: create migrate
create:
	goose -dir=$(MIGRATION_DIR) create $(NAME) sql

migrate:
	./migrate.sh

#goose -v -dir ./migrations postgres "host=localhost port=5432 user=user password=password dbname=news sslmode=disable" down-to 20221120203355