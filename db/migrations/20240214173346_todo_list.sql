-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS todo_list(
                                        id SERIAL PRIMARY KEY,
                                        user_id INTEGER NOT NULL,
                                        task VARCHAR NOT NULL,
                                        active BOOLEAN NOT NULL DEFAULT TRUE,
                                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE todo_list;
-- +goose StatementEnd
