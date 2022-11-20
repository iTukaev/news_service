-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS text_gin ON public.news USING gin(vts_indexed);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS text_gin;
-- +goose StatementEnd
