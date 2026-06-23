-- +goose Up
-- +goose StatementBegin

ALTER TABLE public."validators"
    ADD "enode" text;

ALTER TABLE public."validators"
    ADD "ip" text;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'NOT SUPPORTED';
-- +goose StatementEnd 