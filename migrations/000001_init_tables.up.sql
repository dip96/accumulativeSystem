CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    login CHARACTER varying(50) UNIQUE,
    password bytea NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX ON users (login);

create type order_status as enum ('NEW', 'REGISTERED', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE IF NOT EXISTS orders (
     id serial PRIMARY KEY,
     user_id serial,
     order_id serial UNIQUE,
     accrual numeric,
     status order_status,
     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
     FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX ON orders (user_id);
CREATE INDEX ON orders (order_id);




