.PHONY: up
up:
	docker compose -f docker-compose-local.yml up -d

.PHONY: down
down:
	docker compose -f docker-compose-local.yml down