package service

import (
	"context"
	"time"

	"apiserver/pkg/model"
)

const (
	unknownName = "unknown person"
)

func (s *Service) AddPass(ctx context.Context, mac string, passtime time.Time) (pass *model.PassInfo, err error) {
	defer func(start time.Time) {
		s.logger.Info(ctx, "add pass", "timing", time.Since(start).String())
	}(time.Now())

	if passtime.IsZero() {
		passtime = time.Now()
	}

	pass, err = s.storage.AddPass(ctx, mac, passtime)
	if err != nil {
		s.logger.Error(ctx, "error addpass", "error", err.Error())
		return nil, err
	}

	if pass.Name == "" {
		pass.Name = unknownName
	}

	return pass, nil
}
