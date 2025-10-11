build:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs

run:
	docker compose up

generate:
	sqlc generate

user:
	go test ./internal/domain/user -v

migrate:
	PGPASSWORD=password psql -h localhost -p 5434 -U user -d county_db -f migrations/006_create_assessment_table.sql