MIGRATION_DIR:=./migrations
.PHONY: create migrate
create:
	goose -dir=$(MIGRATION_DIR) create $(NAME) sql

migrate:
	./migrate.sh

.PHONY: up, build
up:
	docker-compose -f docker-compose.psql.yaml up -d
	make migrate
	docker-compose -f docker-compose.services.yaml up -d

build:
	docker-compose -f docker-compose.psql.yaml build
	docker-compose -f docker-compose.services.yaml build
