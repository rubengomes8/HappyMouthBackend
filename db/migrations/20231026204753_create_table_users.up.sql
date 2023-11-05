CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(60) NOT NULL UNIQUE,
    passhash VARCHAR(60),
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp
);

