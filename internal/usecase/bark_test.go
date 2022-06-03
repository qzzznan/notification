package usecase

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"notification/internal/entity"
	"testing"
)

type test struct {
	name string
	mock func()
	res  interface{}
	err  error
}

var (
	errEmptyDevice = fmt.Errorf("empty device")
)

func bark(t *testing.T) (*BarkUseCase, *MockBarkRepo) {
	t.Helper()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := NewMockBarkRepo(ctl)
	b := NewBark(repo, nil)
	return b, repo
}

func TestBark(t *testing.T) {
	b, r := bark(t)
	tests := []test{
		{
			name: "empty",
			mock: func() {
				r.EXPECT().Store(context.Background(), &entity.BarkDevice{}).Return(errEmptyDevice)
			},
			res: &entity.BarkDevice{},
			err: errEmptyDevice,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			err := b.Register(context.Background(), &entity.BarkDevice{})
			require.ErrorIs(t, err, errEmptyDevice)
		})
	}
}
