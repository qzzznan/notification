package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"notification/internal/entity"
)

var _ Bark = (*BarkUseCase)(nil)

type BarkUseCase struct {
	repo BarkRepo
	api  BarkWebAPI
}

func NewBark(r BarkRepo, a BarkWebAPI) *BarkUseCase {
	return &BarkUseCase{
		repo: r,
		api:  a,
	}
}

func (uc *BarkUseCase) Push(ctx context.Context, key string, msg *entity.APNsMessage) error {
	d, err := uc.repo.Get(ctx, &entity.BarkDevice{DeviceKey: key})
	if err != nil {
		return err
	}
	msg.DeviceToken = d.DeviceToken
	return uc.api.Push(ctx, msg)
}

func (uc *BarkUseCase) Register(ctx context.Context, device *entity.BarkDevice) error {
	if device.DeviceToken == "" {
		return fmt.Errorf("device token is empty")
	}

	if device.DeviceKey != "" {
		obj, err := uc.repo.Get(ctx, &entity.BarkDevice{DeviceKey: device.DeviceKey})
		if err != nil {
			return err
		}
		if obj.DeviceToken == device.DeviceToken {
			return nil
		}
	}

	obj, err := uc.repo.Get(ctx, &entity.BarkDevice{DeviceToken: device.DeviceToken})
	if errors.Is(err, sql.ErrNoRows) {
		device.DeviceKey = uuid.NewV4().String()
		return uc.repo.Store(ctx, device)
	}
	if err != nil {
		return err
	}
	device.DeviceKey = obj.DeviceKey
	return nil
}
