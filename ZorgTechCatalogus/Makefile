setup:
	go get -u github.com/swaggo/swag/cmd/swag
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g ./cmd/server/main.go -o ./docs
	go get -u github.com/swaggo/gin-swagger
	go get -u github.com/swaggo/files

build-docker:
	docker compose build --no-cache

run-local:
	export POSTGRES_DB=catalogus
	export POSTGRES_HOST=postgresdb.arjanslab.nl
	export POSTGRES_USER=postgres
	export POSTGRES_PASSWORD=bw4f2LYrDbDzcSYk
	export POSTGRES_PORT=5432
	export MONGODB_USERNAME=mongodb
	export MONGODB_PASSWORD=dESTS7HRtyp23rbm
	export MONGODB_HOST=mongodb.arjanslab.nl
	export MONGODB_PORT=27017
	export MONGODB_DB=ZorgTechCatalogus
	export REDIS_HOST=redis.arjanslab.nl
	export REDIS_PORT=6379
	export REDIS_DB=2
	export REDIS_PASSWORD=KDJD37EzSGzWTKyw
	export JWT_SECRET_KEY=ObL89O3nOSSEj6tbdHako0cXtPErzBUfq8l8o/3KD9g=INSECURE
	export API_SECRET_KEY=cJGZ8L1sDcPezjOy1zacPJZxzZxrPObm2Ggs1U0V+fE=INSECURE
      
	go run cmd/server/main.go

up:
	docker compose up

down:
	docker compose down

restart:
	docker compose restart

build:
	go build -v ./...

test:
	go test -v ./... -race -cover

clean:
	docker stop ZorgTechCatalogus
	docker rm ZorgTechCatalogus
	docker image rm ZorgTechCatalogus-backend

