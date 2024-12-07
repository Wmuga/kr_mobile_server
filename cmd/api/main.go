package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"apiserver/keys"
	"apiserver/pkg/logger"
	"apiserver/pkg/model"
	"apiserver/pkg/server"
	"apiserver/pkg/service"
	"apiserver/pkg/storage"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	block, _ := pem.Decode([]byte(keys.PrivateKey))
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error decode key", err)
		os.Exit(1)
	}

	var configFile string

	flag.StringVar(&configFile, "c", "config.json", "path to config file")
	flag.Usage = func() {
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()

	var cfg = &model.Config{}
	file, err := os.Open(configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cant read config from file", configFile, err)
		os.Exit(1)
	}

	err = json.NewDecoder(file).Decode(cfg)
	file.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cant decode config from file", configFile, err)
		os.Exit(1)
	}

	logger, err := logger.New(cfg.DBDriver, cfg.DBConnectionString)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cant init logger", configFile, err)
		os.Exit(1)
	}

	db, err := storage.New(cfg.DBDriver, cfg.DBConnectionString, cfg.DBMaxConnections, cfg.DBBatchSize)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cant init storage", configFile, err)
		os.Exit(1)
	}

	logic := service.New(logger, db)
	serv := server.New(logic, logger, cfg.Port, cfg.LocalhostOnly, cfg.CheckAuth, key)
	errChan := serv.Start()
	fmt.Println("Start server on ", cfg.Port)
	select {
	case err = <-errChan:
		fmt.Println("error", err)
	case <-sigChan:
		fmt.Println("Stopping")
		serv.Stop()
	}
}
