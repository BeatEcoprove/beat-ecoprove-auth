-- Create Authenticate Table
create table auths(
    id uuid primary key,
    email varchar(50) default not null,
    password text default not null,
    salt text default not null,
    is_active boolean default false,
    role int default 0, -- 0 USER / 1 ORGANIZATION / 2 SPONSOR / 3 ADMIN
    primary key (id)
)

-- Create Profile Table
create table profiles(
    id uuid primary key,
    auth_id uuid default not null,
    type int default 0, -- 0 MAIN / 1 SUB
    primary key (id),
    CONSTRAINT auth_id_fk
    FOREIGN KEY (auth_id)
    REFERENCES auths (id)
)
