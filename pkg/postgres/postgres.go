package postgres

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/go-redis/cache/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"notification/pkg/rdb"
)

type Postgres struct {
	X       *sqlx.DB
	Builder squirrel.StatementBuilderType
	r       *rdb.RDB
}

func New(dsn string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{}

	for _, v := range opts {
		v(pg)
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	pg.X = db
	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return pg, nil
}

func (p *Postgres) Close() {
	if p.X != nil {
		_ = p.X.Close()
	}
}

var ErrNotSetCache = errors.New("not set cache")

func (p *Postgres) SetCache(ctx context.Context, key string, obj interface{}) error {
	if p.r == nil {
		return ErrNotSetCache
	}
	return p.r.Cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   p.r.Key(key),
		Value: obj,
		TTL:   -1,
	})
}

func (p *Postgres) GetCache(ctx context.Context, key string, v interface{}) error {
	if p.r == nil {
		return ErrNotSetCache
	}
	return p.r.Cache.Get(ctx, p.r.Key(key), v)
}

func (p *Postgres) AddListCache(ctx context.Context, key, val string) error {
	if p.r == nil {
		return ErrNotSetCache
	}
	var err error
	arr := make([]string, 0)
	if p.r.Cache.Exists(ctx, key) {
		err = p.r.Cache.Get(ctx, key, arr)
		if err != nil {
			return err
		}
		for _, v := range arr {
			if v == val {
				return nil
			}
		}
	}
	arr = append(arr, val)
	return p.SetCache(ctx, key, arr)
}

func (p *Postgres) DelListCache(ctx context.Context, key, val string) error {
	if p.r == nil {
		return ErrNotSetCache
	}
	if !p.r.Cache.Exists(ctx, key) {
		return nil
	}
	arr := make([]string, 0)
	err := p.r.Cache.Get(ctx, key, arr)
	if err != nil {
		return err
	}
	for i, v := range arr {
		if v == v {
			arr = append(arr[:i], arr[i+1:]...)
			break
		}
	}
	return p.SetCache(ctx, key, arr)
}
