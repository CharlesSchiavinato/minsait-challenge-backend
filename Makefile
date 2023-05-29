include config.env

swag-install:
	which swag || go get github.com/swaggo/swag/cmd/swag@latest

swag: swag-install
	swag init -g server.go -o "./" --outputTypes "yaml"

docker-compose-up:
	docker-compose up -d

docker-compose-stop:
	docker-compose stop

docker-build:
	docker build -t minsait-cash:latest .

docker-run:
	docker run \
	-p 9000:9000 \
	--net minsait-challenge-backend_minsait-network \
	--env SERVER_LOG_JSON_FORMAT=false \
	--env DB_DRIVER=postgres \
	--env DB_URL=postgres://userminsait:Minsait@123@minsait-postgres:5432/db_minsait?sslmode=disable \
	--env DB_MIGRATION_URL=file://migration \
	--env CACHE_URL=redis://:@minsait-redis:6379/0 \
	minsait-cash

migrate-up:
	migrate -source ${DB_MIGRATION_URL} -database "${DB_URL}" up
	
migrate-down:
	migrate -source ${DB_MIGRATION_URL} -database "${DB_URL}" down
	
go-test: 
	go test -v -cover ./...

go-run:
	go run server.go

.PHONY: swag docker-compose-up docker-compose-stop docker-build docker-run migrate-up migrate-down go-test go-run