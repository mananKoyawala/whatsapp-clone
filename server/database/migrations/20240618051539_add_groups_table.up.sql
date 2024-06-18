CREATE TABLE groups (
    "id" SERIAL PRIMARY KEY,
    "admin_id" INT NOT NULL,
    "name" VARCHAR NOT NULL,
    "about" VARCHAR NOT NULL,
    "image" VARCHAR NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,
    CONSTRAINT fk_admin FOREIGN KEY (admin_id) REFERENCES users(id)
)

