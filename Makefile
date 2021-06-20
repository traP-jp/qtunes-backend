.PHONY: up
up:
	@docker compose up

.PHONY: up-d
up-d:
	@docker compose up -d

.PHONY: logs
logs:
	@docker compose logs -f

.PHONY: down
down:
	@docker compose down

.PHONY: lint
lint:
	@docker run --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:latest golangci-lint run
