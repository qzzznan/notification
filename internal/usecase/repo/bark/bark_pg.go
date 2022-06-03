package bark

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"notification/internal/entity"
	"notification/internal/usecase"
	"notification/internal/usecase/repo"
	"notification/pkg/postgres"
)

var _ usecase.BarkRepo = (*Repo)(nil)

const DeviceTable = "t_bark_device"

type Repo struct {
	p *postgres.Postgres
}

func New(p *postgres.Postgres) *Repo {
	return &Repo{p}
}

func (r *Repo) Store(ctx context.Context, device *entity.BarkDevice) error {
	if device.DeviceToken == "" ||
		device.DeviceKey == "" {
		return fmt.Errorf("invalid device")
	}
	sql, args, err := r.p.Builder.
		Insert(DeviceTable).
		Columns("device_key, device_token").
		Values(device.DeviceKey, device.DeviceToken).
		ToSql()
	if err != nil {
		return err
	}

	repo.LogSQL(ctx, sql, args...)

	_, err = r.p.X.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	repo.LogCacheError(ctx, r.p.SetCache(ctx, device.DeviceToken, device))
	repo.LogCacheError(ctx, r.p.SetCache(ctx, device.DeviceKey, device))

	return nil
}

func (r *Repo) Get(ctx context.Context, device *entity.BarkDevice) (*entity.BarkDevice, error) {
	b := r.p.Builder.
		Select("device_token", "device_key").
		From(DeviceTable)

	obj := &entity.BarkDevice{}
	if device.DeviceToken != "" {
		err := r.p.GetCache(ctx, device.DeviceToken, obj)
		if err != nil {
			repo.LogCacheError(ctx, err)
		} else {
			return obj, nil
		}
		b = b.Where(squirrel.Eq{"device_token": device.DeviceToken})
	} else if device.DeviceKey != "" {
		err := r.p.GetCache(ctx, device.DeviceKey, obj)
		if err != nil {
			repo.LogCacheError(ctx, err)
		} else {
			return obj, nil
		}
		b = b.Where(squirrel.Eq{"device_key": device.DeviceKey})
	} else {
		return nil, fmt.Errorf("invalid device")
	}

	sql, args, err := b.ToSql()
	if err != nil {
		return nil, err
	}

	repo.LogSQL(ctx, sql, args...)

	err = r.p.X.GetContext(ctx, obj, sql, args...)
	return obj, err
}
