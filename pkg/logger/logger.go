package logger

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/segmentio/ksuid"
)

type Logger struct {
	db *sql.DB
}

func New(driver, connection string) (*Logger, error) {
	db, err := sql.Open(driver, connection)
	if err != nil {
		return nil, err
	}

	return &Logger{db}, nil
}

// Error implements model.Logger.
func (l *Logger) Error(ctx context.Context, msg string, kv ...string) {
	ctx = context.WithoutCancel(ctx)
	go l.insert(ctx, "error", msg, kv...)
}

// Info implements model.Logger.
func (l *Logger) Info(ctx context.Context, msg string, kv ...string) {
	ctx = context.WithoutCancel(ctx)
	go l.insert(ctx, "info", msg, kv...)
}

func (l *Logger) insert(ctx context.Context, level string, msg string, kv ...string) {
	reqIdInterface := ctx.Value("requestid")
	var (
		reqId  string
		timing string
	)
	if reqIdInterface != nil {
		reqId, _ = reqIdInterface.(string)
	}

	uid := ksuid.New().String()

	payload := make(map[string]string, len(kv)/2+3)
	payload["level"] = level
	payload["requst_id"] = reqId
	payload["msg"] = msg

	for i := 0; i < len(kv)/2; i++ {
		key := kv[i*2]
		value := kv[i*2+1]
		payload[key] = value
		if key == "timing" {
			timing = value
		}
	}

	_, err := l.db.ExecContext(ctx, insert, uid, level, reqId, timing, msg, payload)
	if err != nil {
		fmt.Println(err)
	}
}
