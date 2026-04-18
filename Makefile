include .env
export

env-up:
	docker compose up -d

env-down:
	docker compose down

env-cleanup:
	powershell -Command '$$answer = Read-Host "Delete volumes? (y/n)"; if ($$answer -eq "y" -or $$answer -eq "Y") { docker compose down -v } else { Write-Host "Cancelled" }'

# env-port-forward:
# 	docker compose run --rm ODRProcesing-port-forwarder
	
# env-port-close:
# 	docker compose down port-forwarder

migrate-create:
	docker compose run --rm ODRProcesing-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq $(seq)

migrate-up:
	make migrate-action action=up

migrate-down:
	make migrate-action action=down

migrate-action:
	docker compose run --rm \
		--entrypoint migrate \
		ODRProcesing-postgres-migrate \
		-path=/migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@ODRProcesing-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		${action}

todoapp-run:
 go run cmd/processing/main.go