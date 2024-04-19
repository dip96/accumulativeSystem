CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login CHARACTER varying(50) UNIQUE,
    password BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX ON users (login);

create type order_status as enum ('NEW', 'REGISTERED', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INT,
    order_id BIGINT UNIQUE,
    accrual NUMERIC,
    status order_status,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX ON orders (user_id);
CREATE INDEX ON orders (order_id);


CREATE TABLE IF NOT EXISTS balance_transactions (
    id SERIAL PRIMARY KEY,
    user_id INT,
    amount NUMERIC,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX ON balance_transactions (user_id);

CREATE TABLE IF NOT EXISTS user_balances (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE,
    balance NUMERIC DEFAULT 0,
    withdrawn_balance NUMERIC DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX ON user_balances (user_id);

