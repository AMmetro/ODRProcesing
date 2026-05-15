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
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@odrprocesing-postgres:5432/${POSTGRES_DB}?sslmode=disable \
	${action}

logs-cleanup:
	@read -p "Очистить все log файлы? Опасность утери логов. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf $(PROJECT_ROOT)/out/logs && \
		echo "Файлы логов очищены"; \
	else \
		echo "Очистка логов отменена"; \
	fi

odrprocesing-run:
# unix
# 	export LOGGER_FOLDER=$(PROJECT_ROOT)/out/logs && \
# 	go mod tidy && \
# 	go run cmd/processing/main.go
	powershell -Command "$$env:LOGGER_FOLDER='$(PROJECT_ROOT)/out/logs'; go mod tidy; go run cmd/processing/main.go"
	go mod tidy && \
	go run cmd/processing/main.go

odrprocesing-deploy:
	docker compose up -d --build ODRProcesing

swagger-gen:
	docker compose run --rm swagger \
	init \
	-g cmd/processing/main.go \
	-o docs \
	--parseInternal \
	--parseDependency






	