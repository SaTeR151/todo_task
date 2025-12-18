-- +goose Up
-- +goose StatementBegin
CREATE TABLE scheduler (
    uuid UUID NOT NULL,
    date BIGINT NOT NULL,
    title TEXT NOT NULL DEFAULT "",
    comment TEXT,
    repeat TEXT NOT NULL

    CONSTRAINT scheduler_pk PRIMARY KEY (uuid)

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE scheduler;
-- +goose StatementEnd
