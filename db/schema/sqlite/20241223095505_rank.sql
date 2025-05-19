-- +goose Up
-- +goose StatementBegin

DROP TABLE IF EXISTS "unfinalized_blocks";

CREATE TABLE IF NOT EXISTS "unfinalized_blocks"
(
    "root" BLOB NOT NULL,
    "slot" bigint NOT NULL,
    "header_ver" int NOT NULL,
    "header_ssz" BLOB NOT NULL,
    "block_ver" int NOT NULL,
    "block_ssz" BLOB NOT NULL,
    "status" integer NOT NULL DEFAULT 0,
    "fork_id" BIGINT NOT NULL DEFAULT 0,
    rank integer NOT NULL DEFAULT 0,
    CONSTRAINT "unfinalized_blocks_pkey" PRIMARY KEY ("root", rank)
);

CREATE INDEX IF NOT EXISTS "unfinalized_blocks_slot_idx"
    ON "unfinalized_blocks" 
    ("slot" ASC);

-- SQLite does not have drop constraint
CREATE TABLE IF NOT EXISTS new_slots
(
    slot bigint NOT NULL,
    proposer bigint NOT NULL,
    status smallint NOT NULL,
    root BLOB NOT NULL,
    parent_root BLOB NULL,
    state_root BLOB NULL,
    graffiti BLOB NULL,
    graffiti_text text NULL,
    attestation_count integer NULL DEFAULT 0,
    deposit_count integer NULL DEFAULT 0,
    exit_count integer NULL DEFAULT 0,
    withdraw_count integer NULL DEFAULT 0,
    withdraw_amount bigint NULL DEFAULT 0,
    attester_slashing_count integer NULL DEFAULT 0,
    proposer_slashing_count integer NULL DEFAULT 0,
    bls_change_count integer NULL DEFAULT 0,
    eth_transaction_count integer NULL DEFAULT 0,
    eth_block_number bigint NULL,
    eth_block_hash BLOB NULL,
    eth_block_extra BLOB NULL,
    eth_block_extra_text text NULL,
    sync_participation real NULL DEFAULT 0,
    fork_id BIGINT NOT NULL DEFAULT 0,
    rank integer NOT NULL DEFAULT 0,
    CONSTRAINT slots_pkey PRIMARY KEY (slot, root, rank)
);

CREATE INDEX IF NOT EXISTS "slots_root_idx"
    ON "new_slots"
    ("root" ASC);

CREATE INDEX IF NOT EXISTS "slots_graffiti_idx"
    ON "new_slots"  
    ("graffiti_text" ASC);

CREATE INDEX IF NOT EXISTS "slots_eth_block_extra_idx"
    ON "new_slots"  
    ("eth_block_extra_text" ASC);

CREATE INDEX IF NOT EXISTS "slots_slot_idx"
    ON "new_slots" 
    ("slot" ASC);

CREATE INDEX IF NOT EXISTS "slots_state_root_idx"
    ON "new_slots" 
    ("state_root" ASC );

CREATE INDEX IF NOT EXISTS "slots_eth_block_number_idx"
    ON "new_slots" 
    ("eth_block_number" ASC );

CREATE INDEX IF NOT EXISTS "slots_eth_block_hash_idx"
    ON "new_slots" 
    ("eth_block_hash" ASC );

CREATE INDEX IF NOT EXISTS "slots_proposer_idx"
    ON "new_slots" 
    ("proposer" ASC );

CREATE INDEX IF NOT EXISTS "slots_status_idx"
    ON "new_slots" 
    ("status" ASC );

CREATE INDEX IF NOT EXISTS "slots_parent_root_idx"
    ON "new_slots" 
    ("parent_root" ASC );

CREATE INDEX IF NOT EXISTS "slots_fork_id_idx"
    ON "slots" 
    ("fork_id" ASC);

-- migrate blocks

INSERT INTO "new_slots" (
    slot, proposer, status, root, parent_root, state_root, graffiti, graffiti_text,
    attestation_count, deposit_count, exit_count, withdraw_count, withdraw_amount, attester_slashing_count, 
    proposer_slashing_count, bls_change_count, eth_transaction_count, eth_block_number, eth_block_hash, sync_participation, fork_id
)
SELECT 
    slot, proposer, status, root, parent_root, state_root, graffiti, graffiti_text,
    attestation_count, deposit_count, exit_count, withdraw_count, withdraw_amount, attester_slashing_count, 
    proposer_slashing_count, bls_change_count, eth_transaction_count, eth_block_number, eth_block_hash, sync_participation, fork_id
FROM slots;

DROP TABLE IF EXISTS "slots";

ALTER TABLE "new_slots" RENAME TO "slots"

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'NOT SUPPORTED';
-- +goose StatementEnd