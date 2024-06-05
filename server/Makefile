postgresinit:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres:15-alpine

postgres:
	docker exec -it postgres15 psql -U postgres 

createdb:
	docker exec -it postgres15 createdb --username=postgres --owner=postgres whatsapp

dropdb:
	docker exec -it postgres15 dropdb whatsapp

# migrate create -ext sql -dir database/migrations add_user_table
createtables:
	docker exec -it postgres15 psql -U postgres
	CREATE TABLE users ( "id" bigserial PRIMARY KEY,"name" varchar NOT NULL,"mobile" varchar NOT NULL,"about" varchar NOT NULL,"image" varchar NOT NULL,"last_seen" TIMESTAMP NOT NULL,"is_online" BOOLEAN NOT NULL,"token" varchar NOT NULL,"refresh_token" varchar NOT NULL,"created_at" TIMESTAMP NOT NULL,"updated_at" TIMESTAMP NOT NULL);

	CREATE TABLE otps ("otp" VARCHAR NOT NULL,"expires_at"  TIMESTAMP NOT NULL,"id" int NOT NULL);

migrateup:
	migrate -path database/migrations -database "postgresql://postgres:password@localhost:5432/whatsapp?sslmode=disable" -verbose up

migratedown:
	migrate -path database/migrations -database "postgresql://postgres:password@localhost:5432/whatsapp?sslmode=disable" -verbose down

.PHONY: posgresinit postgres createdb dropdb migrateup migratedown