-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY,
	balance DOUBLE PRECISION
);

CREATE TABLE IF NOT EXISTS transactions (
	id INTEGER PRIMARY KEY generated always as identity,
	sender_id INTEGER REFERENCES users (id),
	money DOUBLE PRECISION, 
	reciever_id INTEGER REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS history (
	id INTEGER PRIMARY KEY generated always as identity,
	user_id INTEGER REFERENCES users (id),
	operation_type VARCHAR(10),
	amount DOUBLE PRECISION,
	time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
