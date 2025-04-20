-- +goose Up
-- +goose StatementBegin
DELETE FROM permission;
INSERT INTO permission (permission_id, resource_name, min_role_priority)
VALUES
    (1, '/user_v1.userV1/Create', 100),
    (3, '/user_v1.userV1/Delete', 100),
    (5, '/chat_v1.chatV1/Delete', 100);


DO $$
BEGIN
    IF EXISTS (
        SELECT FROM information_schema.columns
        WHERE table_name = 'auth_user' AND column_name = 'role'
    ) THEN
        ALTER TABLE auth_user RENAME COLUMN role TO role_id;
    END IF;

    -- Если таблица не существует, создаем её
    IF NOT EXISTS (
        SELECT FROM information_schema.tables
        WHERE table_name = 'auth_user'
    ) THEN
        CREATE TABLE IF NOT EXISTS auth_user (
            id serial PRIMARY KEY,
            name text NOT NULL,
            email text NOT NULL,
            password text NOT NULL,
            role_id INT NOT NULL,
            FOREIGN KEY (role_id) REFERENCES user_role(role_id),
            created_at timestamp NOT NULL DEFAULT now(),
            updated_at timestamp
        );
    ELSIF NOT EXISTS (
        SELECT FROM information_schema.columns
        WHERE table_name = 'auth_user' AND column_name = 'role_id'
    ) THEN
        ALTER TABLE auth_user ADD COLUMN role_id INT NOT NULL;
        ALTER TABLE auth_user ADD CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES user_role (role_id);
    END IF;
END
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
