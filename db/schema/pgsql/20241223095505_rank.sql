-- +goose Up
-- +goose StatementBegin

ALTER TABLE IF EXISTS public.slots
    ADD COLUMN rank integer NOT NULL DEFAULT 0;

ALTER TABLE IF EXISTS public.unfinalized_blocks
    ADD COLUMN rank integer NOT NULL DEFAULT 0;

ALTER TABLE IF EXISTS public.unfinalized_blocks DROP CONSTRAINT IF EXISTS unfinalized_blocks_pkey;

ALTER TABLE IF EXISTS public.unfinalized_blocks
    ADD PRIMARY KEY (root, rank);

ALTER TABLE IF EXISTS public.slots DROP CONSTRAINT IF EXISTS slots_pkey;

ALTER TABLE IF EXISTS public.slots
   ADD PRIMARY KEY (slot, root, rank);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'NOT SUPPORTED';
-- +goose StatementEnd