package postgres

const (
	newsTable = "news"

	idField          = "id"
	titleField       = "title"
	linkField        = "link"
	pubDateField     = "pub_date"
	descriptionField = "description"
	articleField     = "article"
	vtsIndexedField  = "vts_indexed"
)

const (
	insertNews = `INSERT INTO public.news(id, title, link, pub_date, description, article, vts_indexed)
VALUES ($1, $2, $3, $4, $5, $6, to_tsvector(coalesce($5,'') || ' ' || coalesce($6, '')))
`

	getNews = `SELECT title, link, pub_date, description, article
FROM public.news
WHERE vts_indexed @@ phraseto_tsquery($1)
ORDER BY ts_rank_cd(vts_indexed, phraseto_tsquery($1)) DESC
LIMIT 1
`

	getList = `SELECT title, link, pub_date, description, article
FROM public.news
ORDER BY pub_date DESC
LIMIT $1
OFFSET $2`
)
