-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.news (
    id          int8 NOT NULL UNIQUE PRIMARY KEY,
    title       varchar(200) NOT NULL,
    link        varchar(200) NOT NULL,
    pub_date    timestamptz NOT NULL,
    description text NOT NULL,
    article     text,
    vts_indexed tsvector
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.news;
-- +goose StatementEnd
