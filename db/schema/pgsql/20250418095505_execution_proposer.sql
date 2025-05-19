-- +goose Up
-- +goose StatementBegin

ALTER TABLE IF EXISTS public."unfinalized_execution_blocks"
    ADD COLUMN proposer integer;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'NOT SUPPORTED';
-- +goose StatementEnd

