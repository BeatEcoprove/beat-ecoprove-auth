-- +goose Up
-- +goose StatementBegin
create table profiles(
    id uuid not null,
    auth_id uuid not null,
    role int default 0, -- 0 MAIN / 1 SUB
    created_at timestamp default now(),
    updated_at timestamp default now(),
    deleted_at timestamp default null,
    primary key (id),
    CONSTRAINT auth_id_fk
    FOREIGN KEY (auth_id)
    REFERENCES auths (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table profiles;
-- +goose StatementEnd
