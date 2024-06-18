CREATE TABLE group_members (
    "id" SERIAL PRIMARY KEY,
    "g_id" INT NOT NULL,
    "u_id" INT NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,
    CONSTRAINT fk_gid FOREIGN KEY (g_id) REFERENCES groups(id),    
    CONSTRAINT fk_uid FOREIGN KEY (u_id) REFERENCES users(id)    
)