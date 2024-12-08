package logger

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/segmentio/ksuid"
)

const (
	pingTimeout = time.Second * 3
)

type Logger struct {
	db *sql.DB
}

func New(ctx context.Context, driver, connection string) (*Logger, error) {
	db, err := sql.Open(driver, connection)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()
	err = db.PingContext(ctx)

	return &Logger{db}, err
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

	payload := make(payload, len(kv)/2+3)
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
