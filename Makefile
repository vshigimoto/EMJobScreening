# ==============================================================================
# Main
runPerson:
	echo "Starting person service"
	go run ./cmd/person/main.go

# ==============================================================================
# Go migrate postgresql
migrate_up:
	goose -dir migration/person postgres "host=localhost user=postgres password=postgres port=5432 dbname=person sslmode=disable" up 20240122072610_add_user_table.sql

migrate_down:
	goose -dir migration/person postgres "host=localhost user=postgres password=postgres port=5432 dbname=person sslmode=disable" down 20240122072610_add_user_table.sql

# ==============================================================================
# Tools commands
fix-lint:
	golangci-lint run ./...

# ==============================================================================
# Docker compose commands
local:
	echo "Starting local environment"
	docker-compose -f docker/docker-compose.yml up --build