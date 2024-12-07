package service

import (
	"context"
	"time"

	"apiserver/pkg/model"
)

type Storage interface {
	AddPass(ctx context.Context, mac string, time time.Time) (pass *model.PassInfo, err error)
	SelectAll(ctx context.Context, offset int) ([]*model.PassInfo, error)
	SelectToday(ctx context.Context, offset int) ([]*model.PassInfo, error)
}

type Service struct {
	logger  model.Logger
	storage Storage
}

func New(logger model.Logger, storage Storage) *Service {
	return &Service{
		logger:  logger,
		storage: storage,
	}
}
