-- +goose Up
-- +goose StatementBegin
create table auths(
    id uuid not null,
    email varchar(50) not null,
    password text not null,
    salt text not null,
    is_active boolean default false,
    role int default 0, -- 0 USER / 1 ORGANIZATION / 2 SPONSOR / 3 ADMIN
    created_at timestamp default now(),
    updated_at timestamp default now(),
    deleted_at timestamp default null,
    primary key (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table auths;
-- +goose StatementEnd
