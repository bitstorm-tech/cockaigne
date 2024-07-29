-- +goose Up
-- +goose StatementBegin
alter table deals
add column start_instantly boolean not null default false,
add column own_end_date boolean not null default false;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table deals
drop column start_instantly,
drop column own_end_date;
-- +goose StatementEnd
