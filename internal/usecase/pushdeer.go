package usecase

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"notification/internal/entity"
	"notification/pkg/wsManager"
	"strconv"
	"time"
)

var _ PushDeer = (*PushDeerUseCase)(nil)

type PushDeerUseCase struct {
	repo PushDeerRepo
	api  PushDeerWebAPI
	ws   *wsManager.WsManager
}

func NewPushDeer(r PushDeerRepo, a PushDeerWebAPI) *PushDeerUseCase {
	return &PushDeerUseCase{
		repo: r,
		api:  a,
		ws:   wsManager.New(),
	}
}

func (p *PushDeerUseCase) ValidateToken(ctx context.Context, token string) error {
	return nil
}

func (p *PushDeerUseCase) Register(ctx context.Context, appleID, email, name string) (string, error) {
	uid := uuid.NewV4().String()
	err := p.repo.StoreUser(ctx, appleID, email, name, uid)
	return uid, err
}

func (p *PushDeerUseCase) GetUser(ctx context.Context, token string) (*entity.User, error) {
	return p.repo.GetUser(ctx, token, "")
}

func (p *PushDeerUseCase) RegisterDevice(ctx context.Context, info *entity.RegInfo) ([]*entity.Device, error) {
	u, err := p.GetUser(ctx, info.Token)
	if err != nil {
		return nil, err
	}
	if info.Type == "" {
		info.Type = "ios"
	}
	uid := strconv.FormatInt(u.ID, 10)
	obj := &entity.Device{
		UserID:   uid,
		DeviceID: info.DeviceID,
		Type:     info.Type,
		IsClip:   info.IsClip,
		Name:     info.Name,
	}
	err = p.repo.StoreDevice(ctx, obj)
	if err != nil {
		return nil, err
	}
	return p.repo.GetAllDevice(ctx, uid)
}

func (p *PushDeerUseCase) GetAllDevice(ctx context.Context, token string) ([]*entity.Device, error) {
	u, err := p.repo.GetUser(ctx, token, "")
	if err != nil {
		return nil, err
	}
	uid := strconv.FormatInt(u.ID, 10)
	return p.repo.GetAllDevice(ctx, uid)
}

func (p *PushDeerUseCase) RenameDevice(ctx context.Context, id, name string) error {
	return p.repo.UpdateDeviceName(ctx, id, name)
}

func (p *PushDeerUseCase) RemoveDevice(ctx context.Context, id string) error {
	return p.repo.RemoveDevice(ctx, id)
}

func (p *PushDeerUseCase) PushMessage(ctx context.Context, key string, msg *entity.Message) error {
	keyInfo, err := p.repo.GetPushKey(ctx, 0, "", key)
	if err != nil {
		return err
	}

	msg.UserID = keyInfo.UserID
	msg.PushKeyName = keyInfo.Name
	msg.SendAt = time.Now()

	devices, err := p.repo.GetAllDevice(ctx, keyInfo.UserID)
	if err != nil {
		return err
	}
	return p.api.Push(ctx, devices, msg)
}

func (p *PushDeerUseCase) GetMessage(ctx context.Context, token string, offset, limit uint64) ([]*entity.Message, error) {
	uid, err := p.repo.GetUserID(ctx, token)
	if err != nil {
		return nil, err
	}
	return p.repo.GetMessage(ctx, uid, offset, limit)
}

func (p *PushDeerUseCase) RemoveMessage(ctx context.Context, token, msgID string) error {
	return p.repo.RemoveMessage(ctx, msgID)
}

func (p *PushDeerUseCase) GenNewKey(ctx context.Context, token, name string) ([]*entity.PushKey, error) {
	uid, err := p.repo.GetUserID(ctx, token)
	if err != nil {
		return nil, err
	}
	if name == "" {
		name = gofakeit.Name()
	}

	k := &entity.PushKey{
		UserID: uid,
		Key:    gofakeit.LetterN(32),
		Name:   name,
	}

	err = p.repo.StorePushKey(ctx, k)
	if err != nil {
		return nil, err
	}
	return p.repo.GetAllPushKey(ctx, uid)
}

func (p *PushDeerUseCase) RenameKey(ctx context.Context, token, id, name string) error {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	return p.repo.UpdatePushKey(ctx, &entity.PushKey{
		ID:   i,
		Name: name,
	})
}

func (p *PushDeerUseCase) RegenKey(ctx context.Context, token, id string) error {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	return p.repo.UpdatePushKey(ctx, &entity.PushKey{
		ID:  i,
		Key: gofakeit.LetterN(32),
	})
}

func (p *PushDeerUseCase) GetAllKey(ctx context.Context, token string) ([]*entity.PushKey, error) {
	uid, err := p.repo.GetUserID(ctx, token)
	if err != nil {
		return nil, err
	}
	return p.repo.GetAllPushKey(ctx, uid)
}

func (p *PushDeerUseCase) RemoveKey(ctx context.Context, token, id string) error {
	return p.repo.RemovePushKey(ctx, id)
}

func (p *PushDeerUseCase) Upgrade(ctx context.Context, deviceId string, w http.ResponseWriter, r *http.Request, rspH http.Header) error {
	device, err := p.repo.GetDevice(ctx, deviceId)
	if err != nil {
		return err
	}
	if device.Type == "ios" {
		return errors.New("can't upgrade ios device")
	}
	c, err := p.ws.Upgrade(w, r, rspH)
	if err != nil {
		return err
	}
	p.api.Register(ctx, device.DeviceID, c)
	return nil
}
