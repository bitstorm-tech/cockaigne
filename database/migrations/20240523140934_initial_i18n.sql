-- +goose Up
-- +goose StatementBegin
create table
    i18n
(
    key text primary key,
    de  text not null,
    en  text not null
);

alter table accounts
    add column language text not null default 'de';
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
drop table i18n;

alter table accounts
    drop column if exists language;
-- +goose StatementEnd
