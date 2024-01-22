-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "person"
(
    "id"             serial PRIMARY KEY,
    "name"           VARCHAR(50)         NOT NULL,
    "surname"        VARCHAR(255)        NOT NULL,
    "patronymic"     VARCHAR(255)        NOT NULL,
    "age"            int                 NOT NULL,
    "gender"         VARCHAR(255)        NOT NULL,
    "nationality"    VARCHAR(255)        NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE person;
-- +goose StatementEnd