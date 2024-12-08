package main

import (
	"context"
	"crypto/rsa"
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
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	block, _ := pem.Decode([]byte(keys.PrivateKey))
	keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error decode key", err)
		os.Exit(1)
	}

	key, ok := keyInterface.(*rsa.PrivateKey)
	if !ok {
		fmt.Fprintf(os.Stderr, "Error wrong key type %t\n", keyInterface)
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

	fmt.Println("Parsed config")

	logger, err := logger.New(ctx, cfg.DBDriver, cfg.DBConnectionString)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cant init logger", err)
		os.Exit(1)
	}

	fmt.Println("Init logger")

	db, err := storage.New(ctx, cfg.DBDriver, cfg.DBConnectionString, cfg.DBMaxConnections, cfg.DBBatchSize)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cant init storage", err)
		os.Exit(1)
	}

	fmt.Println("Init storage")

	logic := service.New(logger, db)
	serv := server.New(logic, logger, cfg.Port, cfg.LocalhostOnly, cfg.CheckAuth, key)
	errChan := serv.Start()

	fmt.Println("Start server on ", cfg.Port)
	select {
	case err = <-errChan:
		fmt.Println("error", err)
	case <-ctx.Done():
		fmt.Println("Stopping")
		serv.Stop()
	}
}
