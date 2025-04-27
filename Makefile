# commadn to add all packages
# make all

init:
	@echo "Installing all packages"
	@go get -u github.com/gorilla/mux
	@go get -u github.com/go-sql-driver/mysql
	@go get -u github.com/joho/godotenv

run:
	@echo "Running the server"
	@go run cmd/main.go

test:
	@echo "Running the tests"
	@go test -v ./...

tidy:
	@echo "Tidying up the modules"
	@go mod tidy

migration:
	@echo "Running the migrations"
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@echo "Running the migrations up"
	@go run cmd/migrate/main.go up

migrate-down:
	@echo "Running the migrations down"
	@go run cmd/migrate/main.go down

migrate-fix:
	@echo "Fixing the migrations"
	@go run cmd/migrate/main.go fix

migrate-force:
	@echo "Forcing migration version"
	@go run cmd/migrate/main.go force $(VERSION)

migrate-status:
	@echo "Checking migration status"
	@go run cmd/migrate/main.go status
