## Authboss simple example

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name CHARACTER VARYING NOT NULL,
    password CHARACTER VARYING NOT NULL,
    role NUMERIC NOT NULL
)
WITH (
    OIDS=FALSE
);

CREATE UNIQUE INDEX name_idx ON users (name);

INSERT INTO users(name, password, role) VALUES('admin', '$2a$10$CNbDxEMXMm3Idj9KNKJgCO68z2WEqTRGgDyMD1rQjubsaA67grhpu', 1); /* superpassword */
INSERT INTO users(name, password, role) VALUES('user', '$2a$10$TeK4UjQ92fswye5oo5JLwOZSBThjSqJZac0p3D5kGWdPFHpOgai7K', 2); /* 123456 */
