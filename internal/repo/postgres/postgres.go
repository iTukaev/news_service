package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	errorsPkg "github.com/iTukaev/news_service/internal/customerrors"
	"github.com/iTukaev/news_service/internal/repo/models"
)

type PgxPool interface {
	pgxtype.Querier
	Close()
}

func New(pool *pgxpool.Pool, logger *zap.SugaredLogger) *Repo {
	return &Repo{
		pool:   pool,
		logger: logger,
	}
}

func NewPostgres(ctx context.Context, cfg models.Config, logger *zap.SugaredLogger) (*pgxpool.Pool, error) {
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	logger.Debugln("Postgres connection", psqlConn)

	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		return nil, fmt.Errorf("can't connect to database: %v\n", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database error: %v\n", err)
	}
	return pool, nil
}

type Repo struct {
	pool   PgxPool
	logger *zap.SugaredLogger
}

func (r *Repo) NewsExists(ctx context.Context, id uint32) error {
	query, args, err := squirrel.Select(idField).
		From(newsTable).
		Where(squirrel.Eq{idField: id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "postgres NewsInsert: to sql")
	}
	r.logger.Debugln("NewsInsert", query, args)

	var res uint32
	_ = r.pool.QueryRow(ctx, query, args...).Scan(&res)
	if err != nil {
		r.logger.Errorln("postgres news exists check", err)
		return errorsPkg.ErrUnexpected
	}

	if res != id {
		return nil
	}

	return errorsPkg.ErrNewsAlreadyExists
}

func (r *Repo) NewsInsert(ctx context.Context, news *models.News) error {
	args := []interface{}{
		news.ID,
		news.Title,
		news.Link,
		news.PubDate,
		news.Description,
		news.Article,
	}
	r.logger.Debugln("NewsInsert", insertNews, args)

	if _, err := r.pool.Exec(ctx, insertNews, args...); err != nil {
		r.logger.Errorln("postgres news insert", err)
		return errorsPkg.ErrUnexpected
	}

	return nil
}

func (r *Repo) NewsGet(ctx context.Context, search string) (*models.News, error) {
	args := []interface{}{
		search,
	}

	r.logger.Debugln("UserGet", getNews, args)

	row := r.pool.QueryRow(ctx, getNews, args...)

	var news models.News
	if err := row.Scan(&news.Title, &news.Link, &news.PubDate, &news.Description, &news.Article); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorsPkg.ErrNewsNotFound
		}
		return nil, errors.Wrap(err, "postgres news get: row scan")
	}
	r.logger.Debugln("NewsGet", news.Title)

	return &news, nil
}

func (r *Repo) NewsList(ctx context.Context, params models.ListParams) ([]models.News, error) {
	args := []interface{}{
		params.Limit,
		params.Offset,
	}

	r.logger.Debugln("NewsList", getList, args)

	rows, err := r.pool.Query(ctx, getList, args...)
	if err != nil {
		return nil, errors.Wrap(err, "postgres NewsList: query")
	}

	news := make([]models.News, 0)
	for rows.Next() {
		var n models.News
		if err = rows.Scan(&n.Title, &n.Link, &n.PubDate, &n.Description, &n.Article); err != nil {
			return nil, errors.Wrap(err, "postgres NewsList: row scan")
		}
		news = append(news, n)
	}
	r.logger.Debugln("NewsList", news)

	return news, nil
}

func (r *Repo) Close() {
	r.pool.Close()
	r.logger.Infoln("PostgreSQL connection closed")
}
