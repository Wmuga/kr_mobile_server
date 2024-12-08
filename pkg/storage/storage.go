package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"apiserver/pkg/model"

	_ "github.com/lib/pq"
	"github.com/segmentio/ksuid"
)

const (
	pingTimeout = time.Second * 3
)

type Database struct {
	db        *sql.DB
	batchSize int
}

// New - создать новое подключение к базе
// driver - драйвер подключения
// connection - строка подключения
// maxConnections - максимальное число открытых соединений
// batchSize - максимальный размер пачки
func New(ctx context.Context, driver, connection string, maxConnections, batchSize int) (*Database, error) {
	db, err := sql.Open(driver, connection)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxConnections)

	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()
	err = db.PingContext(ctx)

	return &Database{
		db:        db,
		batchSize: batchSize,
	}, err
}

func (d *Database) AddPass(ctx context.Context, mac string, time time.Time) (pass *model.PassInfo, err error) {
	var (
		uid = ksuid.New().String()
	)

	_, err = d.db.ExecContext(ctx, sqlAddPass, uid, time, mac)
	if err != nil {
		return nil, fmt.Errorf("error add new pass: %w", err)
	}

	res := d.db.QueryRowContext(ctx, sqlSelectResult, uid)
	pass, err = scanPassInfo(res)

	if err != nil {
		return nil, fmt.Errorf("error select new pass: %w", err)
	}

	return pass, nil
}

func (d *Database) SelectAll(ctx context.Context, offset int) ([]*model.PassInfo, error) {
	sql := setLimitOffset(sqlSelectAll, d.batchSize, offset)
	return d.scanArray(ctx, sql)
}

func (d *Database) SelectToday(ctx context.Context, offset int) ([]*model.PassInfo, error) {
	sql := setLimitOffset(sqlSelectToday, d.batchSize, offset)
	return d.scanArray(ctx, sql, time.Now().Truncate(time.Hour*24))
}
