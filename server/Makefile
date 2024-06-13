postgresinit:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres:15-alpine

postgresuninit:
	docker stop postgres15
	docker rm postgres15

postgres:
	docker exec -it postgres15 psql -U postgres 

createdb:
	docker exec -it postgres15 createdb --username=postgres --owner=postgres whatsapp

dropdb:
	docker exec -it postgres15 dropdb --username=postgres whatsapp

createmigration:
	migrate create -ext sql -dir database/migrations your_table_name

migrateup:
	migrate -path database/migrations -database "postgresql://postgres:password@localhost:5432/whatsapp?sslmode=disable" -verbose up

migratedown:
	migrate -path database/migrations -database "postgresql://postgres:password@localhost:5432/whatsapp?sslmode=disable" -verbose down

.PHONY: posgresinit postgres createdb dropdb createmigration migrateup migratedown