package storage

import (
	"apiserver/pkg/model"
	"context"
	"fmt"
	"strconv"
	"strings"
)

type scanner interface {
	Scan(dest ...any) error
}

// scanPassInfo взять текущий PassInfo из сканнера
func scanPassInfo(scanner scanner) (*model.PassInfo, error) {
	res := model.PassInfo{}
	err := scanner.Scan(&res.Name, &res.Position, &res.Name, &res.Position, &res.Mac, &res.PassTime)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func setLimitOffset(sql string, limit, offset int) string {
	sql = strings.Replace(sql, placeholderLIMIT, strconv.Itoa(limit), 1)
	return strings.Replace(sql, placeholderOFFSET, strconv.Itoa(offset), 1)
}

func (d *Database) scanArray(ctx context.Context, sql string, params ...interface{}) ([]*model.PassInfo, error) {
	rows, err := d.db.QueryContext(ctx, sql, params...)
	if err != nil {
		return nil, fmt.Errorf("error query all: %w", err)
	}

	data := make([]*model.PassInfo, 0)

	for rows.Next() {
		info, err := scanPassInfo(rows)
		if err != nil {
			return nil, fmt.Errorf("error scan info: %w", err)
		}
		data = append(data, info)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error rows: %w", err)
	}

	return data, nil
}
