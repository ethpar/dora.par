-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS "unfinalized_execution_blocks" (
    "slot" BIGINT NOT NULL,
    "rank" INT NOT NULL,
    "root" bytea NOT NULL,
    eth_block_number BIGINT,
    eth_block_hash bytea,
    "eth_transaction_count" INT NOT NULL DEFAULT 0,
    "block" bytea,
    CONSTRAINT unfinalized_execution_blocks_pkey PRIMARY KEY (slot, rank)        
);

CREATE INDEX IF NOT EXISTS "unfinalized_execution_blocks_root_idx"
    ON "unfinalized_execution_blocks"
    ("root" ASC);  

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'NOT SUPPORTED';
-- +goose StatementEnd