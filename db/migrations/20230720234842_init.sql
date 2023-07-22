-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS public."explorer_state"
(
    "key" character varying(150) COLLATE pg_catalog."default" NOT NULL,
    "value" text COLLATE pg_catalog."default",
    CONSTRAINT "explorer_state_pkey" PRIMARY KEY ("key")
);

CREATE TABLE IF NOT EXISTS public."blocks"
(
    "root" bytea NOT NULL,
    "slot" bigint NOT NULL,
    "parent_root" bytea NOT NULL,
    "state_root" bytea NOT NULL,
    "orphaned" boolean NOT NULL,
    "proposer" bigint NOT NULL,
    "graffiti" bytea NOT NULL,
    "attestation_count" integer NOT NULL DEFAULT 0,
    "deposit_count" integer NOT NULL DEFAULT 0,
    "exit_count" integer NOT NULL DEFAULT 0,
    "attester_slashing_count" integer NOT NULL DEFAULT 0,
    "proposer_slashing_count" integer NOT NULL DEFAULT 0,
    "bls_change_count" integer NOT NULL DEFAULT 0,
    "eth_transaction_count" integer NOT NULL DEFAULT 0,
    "eth_block_number" bigint NOT NULL DEFAULT 0,
    "eth_block_hash" bytea NOT NULL,
    "sync_participation" real NOT NULL DEFAULT 0,
    CONSTRAINT "Blocks_pkey" PRIMARY KEY ("root")
);

CREATE INDEX IF NOT EXISTS "blocks_graffiti_idx"
    ON public."blocks" 
    ("graffiti" ASC NULLS LAST);

CREATE INDEX IF NOT EXISTS "blocks_slot_idx"
    ON public."blocks" 
    ("root" ASC NULLS LAST);

CREATE TABLE IF NOT EXISTS public."orphaned_blocks"
(
    "root" bytea NOT NULL,
    "header" text COLLATE pg_catalog."default" NOT NULL,
    "block" text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT "orphaned_blocks_pkey" PRIMARY KEY ("root")
);

CREATE TABLE IF NOT EXISTS public."epochs"
(
    "epoch" bigint NOT NULL,
    "validator_count" bigint NOT NULL DEFAULT 0,
    "eligible" bigint NOT NULL DEFAULT 0,
    "voted_target" bigint NOT NULL DEFAULT 0,
    "voted_head" bigint NOT NULL DEFAULT 0,
    "voted_total" bigint NOT NULL DEFAULT 0,
    "block_count" smallint NOT NULL DEFAULT 0,
    "orphaned_count" smallint NOT NULL DEFAULT 0,
    "attestation_count" integer NOT NULL DEFAULT 0,
    "deposit_count" integer NOT NULL DEFAULT 0,
    "exit_count" integer NOT NULL DEFAULT 0,
    "attester_slashing_count" integer NOT NULL DEFAULT 0,
    "proposer_slashing_count" integer NOT NULL DEFAULT 0,
    "bls_change_count" integer NOT NULL DEFAULT 0,
    "eth_transaction_count" integer NOT NULL DEFAULT 0,
    "sync_participation" real NOT NULL DEFAULT 0,
    CONSTRAINT "epochs_pkey" PRIMARY KEY ("epoch")
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'NOT SUPPORTED';
-- +goose StatementEnd
