-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS public.contracts
(
    address text NOT NULL,
    "owner" text NOT NULL,
    is_erc20 boolean,
    name text,
    symbol text,
    created_at timestamp without time zone,
    body text,
    CONSTRAINT contracts_pkey PRIMARY KEY (address)
);


ALTER TABLE IF EXISTS public.contracts
    OWNER to explorer;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'NOT SUPPORTED';
-- +goose StatementEnd

