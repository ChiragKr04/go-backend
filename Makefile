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
