.PHONY: up down dev-backend lint-backend test-backend dev-web lint-web build-web lint

# インフラ
up:
	docker compose up -d

down:
	docker compose down

# バックエンド
dev-backend:
	cd backend && air

lint-backend:
	cd backend && golangci-lint run ./...

test-backend:
	cd backend && go test ./...

# フロントエンド
dev-web:
	cd web && npm run dev

lint-web:
	cd web && npm run lint

build-web:
	cd web && npm run build

# 全体
lint: lint-backend lint-web
