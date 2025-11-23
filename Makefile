DB_DSN = host=localhost port=5432 user=root password=root dbname=prs sslmode=disable
MIGRATIONS_DIR = migrations

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" status

db-connect:
	docker exec -it pr-service-db psql -U root -d prs

db-tables:
	docker exec -it pr-service-db psql -U root -d prs -c "\dt"

db-describe:
	docker exec -it pr-service-db psql -U root -d prs -c "\d users"

up:
	docker-compose up -d

down:
	docker-compose down