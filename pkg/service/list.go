package service

import (
	"apiserver/pkg/model"
	"context"
	"strconv"
	"time"
)

func (s *Service) ListAll(ctx context.Context, offset int) (pass []*model.PassInfo, err error) {
	defer func(start time.Time) {
		s.logger.Info(ctx, "list all pass", "timing", time.Since(start).String(), "offset", strconv.Itoa(offset))
	}(time.Now())

	pass, err = s.storage.SelectAll(ctx, offset)
	if err != nil {
		s.logger.Error(ctx, "error list all pass", "error", err.Error())
		return nil, err
	}

	for i := range pass {
		if pass[i].Name == "" {
			pass[i].Name = unknownName
		}
	}

	return pass, nil
}

func (s *Service) ListToday(ctx context.Context, offset int) (pass []*model.PassInfo, err error) {
	defer func(start time.Time) {
		s.logger.Info(ctx, "list today pass", "timing", time.Since(start).String(), "offset", strconv.Itoa(offset))
	}(time.Now())

	pass, err = s.storage.SelectToday(ctx, offset)
	if err != nil {
		s.logger.Error(ctx, "error list today pass", "error", err.Error())
		return nil, err
	}

	for i := range pass {
		if pass[i].Name == "" {
			pass[i].Name = unknownName
		}
	}

	return pass, nil
}
