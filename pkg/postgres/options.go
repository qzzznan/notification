package postgres

import "notification/pkg/rdb"

type Option func(postgres *Postgres)

func WithCache(r *rdb.RDB) Option {
	return func(postgres *Postgres) {
		postgres.r = r
	}
}
