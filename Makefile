run:
	go run cmd/api/main.go
migrate: 
	go run cmd/migrator/main.go
test: 
	go test -count=1 -v ./internal/db/postgres
docker: 
