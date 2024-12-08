package logger

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type payload map[string]interface{}

func (a payload) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *payload) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
