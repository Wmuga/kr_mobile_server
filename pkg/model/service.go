package model

import (
	"context"
	"encoding/json"
	"time"
)

type Logger interface {
	Info(ctx context.Context, msg string, kv ...string)
	Error(ctx context.Context, msg string, kv ...string)
}

type PassInfo struct {
	Name     string
	Position string
	Mac      string
	PassTime time.Time
}

type passInfo struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	PassTime string `json:"pass_time"`
}

func (p *PassInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(passInfo{
		Name:     p.Name,
		Position: p.Position,
		PassTime: p.PassTime.Format(time.RFC3339),
	})
}
