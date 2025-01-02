-- +goose Up
create table auth_user (
    id serial primary key,
    name text not null,
    email text not null,
    password text not null,
    role text not null,
    created_at timestamp not null default now(),
    updated_at timestamp
);

-- +goose Down
drop table auth_user;
