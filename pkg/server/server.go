package server

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"apiserver/pkg/model"
)

type Service interface {
	AddPass(ctx context.Context, mac string, passtime time.Time) (pass *model.PassInfo, err error)
	ListAll(ctx context.Context, offset int) (pass []*model.PassInfo, err error)
	ListToday(ctx context.Context, offset int) (pass []*model.PassInfo, err error)
}

type Server struct {
	service    Service
	logger     model.Logger
	server     *http.Server
	privateKey *rsa.PrivateKey
}

func New(service Service, logger model.Logger, port int, localhost, checkAuth bool, privateKey *rsa.PrivateKey) *Server {
	listenStr := ":" + strconv.Itoa(port)
	if localhost {
		listenStr = "127.0.0.1" + listenStr
	}

	mux := http.NewServeMux()

	s := &Server{
		service:    service,
		logger:     logger,
		privateKey: privateKey,
	}

	mux.HandleFunc("POST /add", s.add)
	mux.HandleFunc("POST /list/all", s.listAll)
	mux.HandleFunc("POST /list/today", s.listToday)

	var handler http.Handler = mux
	if checkAuth {
		handler = s.authMiddleware(handler)
	}

	server := &http.Server{
		Addr:         listenStr,
		Handler:      s.requestIdMiddleware(s.loggerMiddleware(handler)),
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	s.server = server

	return s
}

func (s *Server) Start() chan error {
	errChan := make(chan error, 1)

	go func() {
		errChan <- s.server.ListenAndServe()
	}()

	return errChan
}

func (s *Server) Stop() {
	s.server.Shutdown(context.Background())
}

func transportError(w http.ResponseWriter, statusCode int, err error) {
	resp := ServiceResponse{
		Status:  statusCode,
		Message: err.Error(),
	}

	bytes, _ := json.Marshal(resp)
	w.WriteHeader(statusCode)
	w.Write(bytes)
}

func transportResponse(w http.ResponseWriter, data interface{}, count int) {
	resp := ServiceResponse{
		Status:  http.StatusOK,
		Message: "OK",
		Count:   count,
		Data:    data,
	}

	bytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
