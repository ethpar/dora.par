-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS public.transactions
(
    "hash"              text NOT NULL,
    "block_number"      bigint,
    "block_rank"        int,
    "created_at"        timestamp,
    nonce               int,
    block_hash          text,
    transaction_index   int,
    "from"              text,
    "to"                text,
    value               numeric,
    gas                 bigint,
    gas_price           bigint,
    is_error            boolean,
    receipt_status      text,
    input               text,
    contract_address    text,
    cumulative_gas_used bigint,
    gas_used            bigint,
    confirmations       int,
    "type" int,
    CONSTRAINT "transactions_pkey" PRIMARY KEY ("hash")
);

CREATE INDEX IF NOT EXISTS "transactions_from_idx"
    ON public."transactions"
        ("from" ASC NULLS LAST);

CREATE INDEX IF NOT EXISTS "transactions_to_idx"
    ON public."transactions"
        ("to" ASC NULLS LAST);

CREATE INDEX IF NOT EXISTS transactions_from_ci_idx ON public.transactions ((lower("from")));

CREATE INDEX IF NOT EXISTS transactions_to_ci_idx ON public.transactions ((lower("to")));

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'NOT SUPPORTED';
-- +goose StatementEnd

