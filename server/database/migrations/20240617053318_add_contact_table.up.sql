CREATE TABLE contacts (
    id bigserial PRIMARY KEY,
    uid INT NOT NULL, 
    cid INT NOT NULL,
    CONSTRAINT fk_uid FOREIGN KEY (uid) REFERENCES users(id),
    CONSTRAINT fk_cid FOREIGN KEY (cid) REFERENCES users(id)
)