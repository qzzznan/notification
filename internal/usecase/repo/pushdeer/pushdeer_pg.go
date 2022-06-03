package pushdeer

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"notification/internal/entity"
	"notification/internal/usecase"
	"notification/internal/usecase/repo"
	"notification/pkg/postgres"
	"strconv"
	"time"
)

var _ usecase.PushDeerRepo = (*Repo)(nil)

const (
	UserTable    = "t_user"
	DeviceTable  = "t_device"
	PushKeyTable = "t_push_key"
	MessageTable = "t_message"
)

type Repo struct {
	p *postgres.Postgres
}

func New(pg *postgres.Postgres) *Repo {
	return &Repo{pg}
}

func (r *Repo) GetUserID(ctx context.Context, token string) (string, error) {
	u, err := r.GetUser(ctx, token, "")
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(u.ID, 10), nil
}

func (r *Repo) StoreUser(ctx context.Context, appleID, email, name, uuid string) error {
	sql, args, err := r.p.Builder.Insert(UserTable).
		Columns("apple_id", "email", "name", "uuid", "created_at").
		Values(appleID, email, name, uuid, time.Now()).
		Suffix("ON CONFLICT (apple_id) DO NOTHING").
		ToSql()

	if err != nil {
		return err
	}

	repo.LogSQL(ctx, sql, args...)

	_, err = r.p.X.ExecContext(ctx, sql, args...)
	return err
}

func (r *Repo) GetUser(ctx context.Context, uuid, appleID string) (*entity.User, error) {
	user := &entity.User{}

	b := r.p.Builder.
		Select("id", "apple_id", "email", "name", "uuid", "created_at").
		From(UserTable)
	if uuid != "" {
		err := r.p.GetCache(ctx, uuid, user)
		if err != nil {
			repo.LogCacheError(ctx, err)
		} else {
			return user, nil
		}
		b = b.Where(squirrel.Eq{"uuid": uuid})
	} else if appleID != "" {
		err := r.p.GetCache(ctx, appleID, user)
		if err != nil {
			repo.LogCacheError(ctx, err)
		} else {
			return user, nil
		}
		b = b.Where(squirrel.Eq{"apple_id": appleID})
	} else {
		return nil, errors.New("uuid or apple_id is required")
	}
	str, args, err := b.ToSql()
	if err != nil {
		return nil, err
	}

	repo.LogSQL(ctx, str, args...)

	err = r.p.X.GetContext(ctx, user, str, args...)
	if err != nil {
		return nil, err
	}

	repo.LogCacheError(ctx, r.p.SetCache(ctx, uuid, user))
	repo.LogCacheError(ctx, r.p.SetCache(ctx, appleID, user))

	return user, nil
}

func (r *Repo) StoreDevice(ctx context.Context, device *entity.Device) error {
	now := time.Now()
	sql, args, err := r.p.Builder.Insert(DeviceTable).
		Columns("user_id", "device_id", "type", "is_clip", "name", "created_at", "updated_at").
		Values(device.UserID, device.DeviceID, device.Type, device.IsClip, device.Name, now, now).
		Suffix("ON CONFLICT (device_id) DO NOTHING").
		ToSql()

	if err != nil {
		return err
	}

	repo.LogSQL(ctx, sql, args...)

	_, err = r.p.X.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	repo.LogCacheError(ctx, r.p.SetCache(ctx, device.DeviceID, device))
	err = r.p.AddListCache(ctx, r.userDeviceCacheKey(device.UserID), device.DeviceID)
	repo.LogCacheError(ctx, err)

	return nil
}

func (r *Repo) userDeviceCacheKey(key string) string {
	return fmt.Sprintf("user_device_%s", key)
}

func (r *Repo) GetDevice(ctx context.Context, deviceID string) (*entity.Device, error) {
	device := &entity.Device{}
	err := r.p.GetCache(ctx, deviceID, device)
	if err != nil {
		repo.LogCacheError(ctx, err)
	}

	sql, args, err := r.p.Builder.Select("device_id", "user_id", "type", "is_clip", "name", "created_at", "updated_at").
		From(DeviceTable).
		Where(squirrel.Eq{"device_id": deviceID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	repo.LogSQL(ctx, sql, args...)

	err = r.p.X.GetContext(ctx, device, sql, args...)
	if err != nil {
		return nil, err
	}

	repo.LogCacheError(ctx, r.p.SetCache(ctx, deviceID, device))
	repo.LogCacheError(ctx, r.p.AddListCache(ctx, r.userDeviceCacheKey(device.UserID), deviceID))

	return device, nil
}

func (r *Repo) GetAllDevice(ctx context.Context, userID string) ([]*entity.Device, error) {
	devices := make([]*entity.Device, 0)
	deviceIDList := make([]string, 0)
	err := r.p.GetCache(ctx, r.userDeviceCacheKey(userID), &deviceIDList)
	if err != nil {
		repo.LogCacheError(ctx, err)
	} else {
		for _, v := range deviceIDList {
			item := new(entity.Device)
			err = r.p.GetCache(ctx, v, item)
			if err != nil {
				repo.LogCacheError(ctx, err)
				break
			}
			devices = append(devices, item)
		}
		return devices, nil
	}

	sql, args, err := r.p.Builder.
		Select("id", "user_id", "device_id", "type", "is_clip", "name", "created_at", "updated_at").
		From(DeviceTable).
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	err = r.p.X.SelectContext(ctx, &devices, sql, args...)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *Repo) UpdateDeviceName(ctx context.Context, deviceID, name string) error {
	sql, args, err := r.p.Builder.Update(DeviceTable).
		Set("name", name).
		Where(squirrel.Eq{"id": deviceID}).
		ToSql()
	if err != nil {
		return err
	}

	repo.LogSQL(ctx, sql, args...)

	_, err = r.p.X.ExecContext(ctx, sql, args...)
	return err
}

func (r *Repo) RemoveDevice(ctx context.Context, s string) error {
	query, args, err := r.p.Builder.Delete(DeviceTable).
		Where(squirrel.Eq{"id": s}).
		ToSql()
	if err != nil {
		return err
	}

	repo.LogSQL(ctx, query, args...)

	_, err = r.p.X.ExecContext(ctx, query, args...)
	return err
}

func (r *Repo) StorePushKey(ctx context.Context, key *entity.PushKey) error {
	now := time.Now()
	query, args, err := r.p.Builder.Insert(PushKeyTable).
		Columns("user_id", "name", "key", "created_at", "updated_at").
		Values(key.UserID, key.Name, key.Key, now, now).
		ToSql()

	if err != nil {
		return err
	}

	repo.LogSQL(ctx, query, args...)

	_, err = r.p.X.ExecContext(ctx, query, args...)
	return err
}

func (r *Repo) GetPushKey(ctx context.Context, id int64, name string, pushKey string) (*entity.PushKey, error) {
	b := r.p.Builder.
		Select("id", "user_id", "name", "key", "created_at", "updated_at").
		From(PushKeyTable)
	if id != 0 {
		b = b.Where(squirrel.Eq{"id": id})
	} else if name != "" {
		b = b.Where(squirrel.Eq{"name": name})
	} else if pushKey != "" {
		b = b.Where(squirrel.Eq{"key": pushKey})
	} else {
		return nil, errors.New("id or name or pushKey is required")
	}

	query, args, err := b.ToSql()
	if err != nil {
		return nil, err
	}

	repo.LogSQL(ctx, query, args...)

	key := &entity.PushKey{}
	err = r.p.X.GetContext(ctx, key, query, args...)
	return key, err
}

func (r *Repo) GetAllPushKey(ctx context.Context, s string) ([]*entity.PushKey, error) {
	query, args, err := r.p.Builder.
		Select("id", "user_id", "name", "key", "created_at", "updated_at").
		From(PushKeyTable).
		Where(squirrel.Eq{"user_id": s}).
		ToSql()

	if err != nil {
		return nil, err
	}

	repo.LogSQL(ctx, query, args...)

	list := make([]*entity.PushKey, 0)
	return list, r.p.X.SelectContext(ctx, &list, query, args...)
}

func (r *Repo) UpdatePushKey(ctx context.Context, key *entity.PushKey) error {
	b := r.p.Builder.Update(PushKeyTable)
	if key.Name != "" {
		b = b.Set("name", key.Name)
	}
	if key.Key != "" {
		b = b.Set("key", key.Key)
	}
	query, args, err := b.Where(squirrel.Eq{"id": key.ID}).ToSql()
	if err != nil {
		return err
	}

	repo.LogSQL(ctx, query, args...)

	_, err = r.p.X.ExecContext(ctx, query, args...)
	return err
}

func (r *Repo) RemovePushKey(ctx context.Context, s string) error {
	query, args, err := r.p.Builder.Delete(PushKeyTable).
		Where(squirrel.Eq{"id": s}).
		ToSql()
	if err != nil {
		return err
	}

	repo.LogSQL(ctx, query, args...)

	_, err = r.p.X.ExecContext(ctx, query, args...)
	return err
}

func (r *Repo) StoreMessage(ctx context.Context, m *entity.Message) error {
	now := time.Now()
	query, args, err := r.p.Builder.Insert(MessageTable).
		Columns("user_id", "text", "type", "note", "push_key", "url", "send_at").
		Values(m.UserID, m.Text, m.Type, m.Note, m.PushKeyName, m.URL, now).
		ToSql()
	if err != nil {
		return err
	}

	repo.LogSQL(ctx, query, args...)

	_, err = r.p.X.ExecContext(ctx, query, args...)
	return err
}

func (r *Repo) GetMessage(ctx context.Context, userID string, offset, count uint64) ([]*entity.Message, error) {
	query, args, err := r.p.Builder.
		Select("id", "user_id", "text", "type", "note", "push_key", "url", "send_at").
		From(MessageTable).
		Where(squirrel.Eq{"user_id": userID}).
		Offset(offset).Limit(count).
		ToSql()
	if err != nil {
		return nil, err
	}

	repo.LogSQL(ctx, query, args...)

	list := make([]*entity.Message, 0)
	return list, r.p.X.SelectContext(ctx, &list, query, args...)
}

func (r *Repo) RemoveMessage(ctx context.Context, s string) error {
	query, args, err := r.p.Builder.Delete(MessageTable).
		Where(squirrel.Eq{"id": s}).
		ToSql()
	if err != nil {
		return err
	}

	repo.LogSQL(ctx, query, args...)

	_, err = r.p.X.ExecContext(ctx, query, args...)
	return err
}
