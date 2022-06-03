package rdb

type Option func(r *RDB)

func WithKeyPrefix(prefix string) Option {
	return func(r *RDB) {
		r.keyPrefix = prefix
	}
}
