CREATE TABLE users ( 
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL,
    "mobile" varchar NOT NULL,
    "about" varchar NOT NULL,
    "image" varchar NOT NULL,
    "last_seen" TIMESTAMP NOT NULL,
    "is_online" BOOLEAN NOT NULL,
    "token" varchar NOT NULL,
    "refresh_token" varchar NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);