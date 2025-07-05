-- +goose Up
-- +goose StatementBegin
ALTER TABLE topics ADD COLUMN status TEXT NOT NULL DEFAULT 'Aguardando Abertura';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE topics DROP COLUMN status;
-- +goose StatementEnd 