postgresinit:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres:15-alpine

postgres:
	docker exec -it postgres15 psql -U postgres 

createdb:
	docker exec -it postgres15 createdb --username=postgres --owner=postgres whatsapp

dropdb:
	docker exec -it postgres15 dropdb whatsapp

# migrate create -ext sql -dir database/migrations add_user_table
migrateup:
	migrate -path database/migrations -database "postgresql://postgres:password@localhost:5432/whatsapp?sslmode=disable" -verbose up

migratedown:
	migrate -path database/migrations -database "postgresql://postgres:password@localhost:5432/whatsapp?sslmode=disable" -verbose down

.PHONY: posgresinit postgres createdb dropdb migrateup migratedown