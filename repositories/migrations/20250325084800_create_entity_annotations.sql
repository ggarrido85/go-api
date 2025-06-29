-- +goose Up

create table entity_annotations (
    id uuid default gen_random_uuid(),
    org_id uuid not null,
    object_type text not null,
    object_id text not null,
    case_id uuid,
    annotation_type text not null,
    payload jsonb not null,
    annotated_by uuid,
    created_at timestamp with time zone not null default now(),
    deleted_at timestamp with time zone default null,

    primary key (id),
    foreign key (org_id) references organizations (id),
    foreign key (case_id) references cases (id) on delete set null,
    foreign key (annotated_by) references users (id) on delete set null
);

create index idx_entity_annotations
    on entity_annotations (org_id, object_type, object_id, annotation_type)
    where deleted_at is null;

create index idx_entity_annotations_case_id
    on entity_annotations (org_id, case_id)
    where deleted_at is null;

-- +goose Down

drop table entity_annotations;
