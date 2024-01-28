CREATE TABLE IF NO EXISTS user_coins (
    user_id SERIAL PRIMARY KEY REFERENCES users(id),
    coins INT NOT NULL,
    created_at timestamp,
    updated_at timestamp,
);