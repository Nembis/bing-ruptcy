
dev-build:
	docker compose -f compose.base.yml -f compose.dev.yml build

dev-up:
	docker compose -f compose.base.yml -f compose.dev.yml up -d

dev-down:
	docker compose -f compose.base.yml -f compose.dev.yml down

dev-logs:
	docker compose -f compose.base.yml -f compose.dev.yml logs -f
