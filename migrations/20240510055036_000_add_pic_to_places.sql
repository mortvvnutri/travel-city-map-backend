-- +goose Up
-- +goose StatementBegin
ALTER TABLE places
ADD pic character varying(128) DEFAULT NULL; 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE places
DROP COLUMN pic; 
-- +goose StatementEnd
