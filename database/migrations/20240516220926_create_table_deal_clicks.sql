-- +goose Up
-- +goose StatementBegin
create table if not exists deal_clicks
(
    deal_id    uuid        not null references deals (id),
    account_id uuid        not null references accounts (id),
    clicked    timestamptz not null default now(),
    constraint "deal_id_account_id_key" unique (deal_id, account_id)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
drop table if exists deal_clicks;
-- +goose StatementEnd
