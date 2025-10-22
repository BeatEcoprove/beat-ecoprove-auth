-- +goose Up
-- +goose StatementBegin
create table auths(
    id uuid not null,
    email varchar(50) not null,
    password text not null,
    salt text not null,
    is_active boolean default false,
    role varchar default 'anonymous',
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
