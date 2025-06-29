-- +goose Up
-- +goose StatementBegin

alter table sanction_check_configs
    add column counterparty_id_expression
    jsonb;

alter table sanction_checks
    add column whitelisted_entities text[] not null default '{}';

alter table sanction_check_matches
    add column counterparty_id text;

create table sanction_check_whitelists (
    id uuid default gen_random_uuid(),
    org_id uuid not null,
    counterparty_id text not null,
    entity_id text not null,
    whitelisted_by uuid not null,
    created_at timestamp with time zone default now(),

    primary key (id),
    constraint fk_organization foreign key (org_id) references organizations (id),
    constraint fk_user foreign key (whitelisted_by) references users (id)
);

create unique index idx_sanction_check_whitelist on sanction_check_whitelists (org_id, counterparty_id, entity_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table sanction_check_configs
    drop column counterparty_id_expression;

alter table sanction_checks
    drop column whitelisted_entities;

alter table sanction_check_matches
    drop column counterparty_id;

drop table sanction_check_whitelists;

-- +goose StatementEnd
