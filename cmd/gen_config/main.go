// Генерация конфига по тегу default используя пакет рефлект
// Пока без вложенных структур
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"apiserver/pkg/model"
)

const outfile = "config-example.json"

func main() {
	cfg := &model.Config{}

	val := reflect.ValueOf(cfg).Elem()
	valType := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		defaultValue := valType.Field(i).Tag.Get("default")
		if defaultValue == "" {
			continue
		}

		err := json.Unmarshal([]byte(defaultValue), field.Addr().Interface())
		if err == nil {
			continue
		}
		err = json.Unmarshal([]byte(`"`+defaultValue+`"`), field.Addr().Interface())
		if err != nil {
			fmt.Fprintln(os.Stderr, field, defaultValue, err)
			os.Exit(1)
		}
	}

	bytes, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = os.WriteFile(outfile, bytes, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
