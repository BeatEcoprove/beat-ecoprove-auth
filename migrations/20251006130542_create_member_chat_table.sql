-- +goose Up
-- +goose StatementBegin
create table member_chat(
    id uuid not null,
    group_id uuid not null,
    member_id uuid not null,
    role varchar default 'member',
    created_at timestamp default now(),
    updated_at timestamp default now(),
    deleted_at timestamp default null,
    primary key (id),
    CONSTRAINT fk_member
        FOREIGN KEY (member_id)
        REFERENCES auths(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table member_chat;
-- +goose StatementEnd
