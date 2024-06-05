DB_CONTAINER_NAME=postgres15
DB_USER=postgres
SQL_FILE=tables.sql


postgresinit:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres:15-alpine

postgres:
	docker exec -it postgres15 psql -U postgres 

createdb:
	docker exec -it postgres15 createdb --username=postgres --owner=postgres whatsapp

dropdb:
	docker exec -it postgres15 dropdb --username=postgres whatsapp

createtables:
	docker cp $(SQL_FILE) $(DB_CONTAINER_NAME):$(SQL_FILE)
	docker exec -i $(DB_CONTAINER_NAME) psql -U $(DB_USER) -f $(SQL_FILE)
	
migrateup:
	migrate -path database/migrations -database "postgresql://postgres:password@localhost:5432/whatsapp?sslmode=disable" -verbose up

migratedown:
	migrate -path database/migrations -database "postgresql://postgres:password@localhost:5432/whatsapp?sslmode=disable" -verbose down

.PHONY: posgresinit postgres createdb dropdb migrateup migratedown