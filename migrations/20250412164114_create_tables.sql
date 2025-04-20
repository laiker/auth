-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS permission (
    permission_id INT PRIMARY KEY,
    resource_name VARCHAR(100) NOT NULL,
    min_role_priority INT NOT NULL DEFAULT 10
);

create table auth_user_log (
    id serial primary key,
    name text not null,
    entity_id int null,
    created_at timestamp not null default now()
);

CREATE TABLE IF NOT EXISTS user_role (
    role_id INT PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL,
    priority INT NOT NULL
);
INSERT INTO user_role (role_id, role_name, priority)
VALUES
    (1, 'user', 10),
    (2, 'admin', 100)
ON CONFLICT (role_id) DO NOTHING;

create table if not exists auth_user (
    id serial primary key,
    name text not null,
    email text not null,
    password text not null,
    role_id INT NOT NULL,
    FOREIGN KEY (role_id) REFERENCES user_role(role_id),
    created_at timestamp not null default now(),
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists permission;
drop table if exists auth_user;
drop table if exists auth_user_log;
drop table if exists user_role;
-- +goose StatementEnd
