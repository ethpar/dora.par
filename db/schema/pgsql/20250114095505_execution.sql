-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS public."unfinalized_execution_blocks"
(
    "slot" bigint NOT NULL,
    "rank" integer NOT NULL,
    "root" bytea NOT NULL,
    eth_block_number bigint,
    eth_block_hash bytea,
    "eth_transaction_count" integer NOT NULL DEFAULT 0,
    "block" bytea,
    CONSTRAINT "unfinalized_execution_blocks_pkey" PRIMARY KEY ("slot", "rank")
    );

CREATE INDEX IF NOT EXISTS "unfinalized_execution_blocks_root_idx"
    ON public."unfinalized_execution_blocks"
        ("root" ASC NULLS LAST);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'NOT SUPPORTED';
-- +goose StatementEnd

