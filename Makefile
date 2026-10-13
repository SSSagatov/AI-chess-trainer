run:
	go run ./cmd/app
migrate: 
	go run ./cmd/migrator/main.go
test: 
	go test -count=1 -v ./internal/db/postgres
docker: 
