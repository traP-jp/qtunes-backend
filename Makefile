.PHONY: up
up:
	@docker compose up

.PHONY: up-d
up-d:
	@docker compose up -d

.PHONY: logs
logs:
	@docker compose logs -f

.PHONY: stop
stop:
	@docker compose stop

.PHONY: down
down:
	@docker compose down

.PHONY: lint
lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --fix

.PHONY: db-dev
db-dev:
	@docker compose exec mariadb mariadb -uroot -ppassword 21hack02
