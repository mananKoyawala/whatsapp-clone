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

CREATE TABLE otps (
    "otp" VARCHAR NOT NULL,
    "expires_at" varchar NOT NULL,
    "id" int NOT NULL
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    message_type TEXT NOT NULL,
    message_text TEXT,
    media_url TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_sender FOREIGN KEY (sender_id) REFERENCES users(id),
    CONSTRAINT fk_receiver FOREIGN KEY (receiver_id) REFERENCES users(id)
);
