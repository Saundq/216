docker-up:
	docker-compose up -d

docker-down:
	docker-compose down --remove-orphans

frontend-install:
	docker-compose run --rm node npm i

run-frontend:
	docker-compose exec node npm run dev -- --port 5173

run-orc:
	cd backend; \
	go run cmd/orchestrator/main.go

run-agent:
	cd backend; \
	go run cmd/agent/main.go


