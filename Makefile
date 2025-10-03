build:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs

run:
	docker compose up

user:
	go test ./internal/domain/user -v


