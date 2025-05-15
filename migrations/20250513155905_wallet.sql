-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS wallets (
     wallet_id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
     balance        BIGINT DEFAULT 0 CHECK (balance >= 0),
     created_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
     updated_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS wallets;