.PHONY:

migrate/up:
	@echo "Running database migrations..."
	@goose -dir sql/migrations postgres postgres://user:password@localhost:5433/mydatabase up
migrate/down:
	@echo "Rolling back database migrations..."
	@goose -dir sql/migrations postgres postgres://user:password@localhost:5433/mydatabase down 

sqlc/generate:
	@echo "Generating SQL code..."
	@sqlc generate

create/migration:
	@echo "Creating new migration..."
	@goose -dir sql/migrations create $(name) sql