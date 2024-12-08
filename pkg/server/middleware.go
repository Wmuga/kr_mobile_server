package server

import (
	"bytes"
	"cmp"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/segmentio/ksuid"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type WriterWrapper struct {
	http.ResponseWriter
	code int
}

func (w *WriterWrapper) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			transportError(w, http.StatusBadRequest, err)
			return
		}

		if len(body) == 0 {
			transportError(w, http.StatusUnauthorized, ErrUnauthorized)
			return
		}

		ok, err := s.checkAuth(body)
		if err != nil {
			transportError(w, http.StatusBadRequest, err)
			return
		}

		if !ok {
			transportError(w, http.StatusUnauthorized, ErrUnauthorized)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		next.ServeHTTP(w, r)
	})
}

func (s *Server) loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &WriterWrapper{ResponseWriter: w, code: http.StatusOK}

		next.ServeHTTP(writer, r)

		s.logger.Info(r.Context(), "query "+r.URL.Path, "method", r.Method, "code", strconv.Itoa(writer.code))
	})
}

func (s *Server) requestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqid := ctx.Value("requestid")
		if _, ok := reqid.(string); ok {
			next.ServeHTTP(w, r)
			return
		}

		ctx = context.WithValue(ctx, "requestid", ksuid.New().String())
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (s *Server) checkAuth(body []byte) (bool, error) {
	type kv struct {
		key   string
		value string
	}

	req := map[string]interface{}{}
	err := json.Unmarshal(body, &req)
	if err != nil {
		return false, err
	}

	tokenInterface, ok := req["token"]
	if !ok {
		return false, nil
	}

	tokenString := fmt.Sprint(tokenInterface)

	kvs := make([]kv, 0, len(req))
	for k, v := range req {
		if k == "token" {
			continue
		}

		kvs = append(kvs, kv{
			key:   k,
			value: fmt.Sprint(v),
		})
	}

	slices.SortFunc(kvs, func(l, r kv) int {
		return cmp.Compare(l.key, r.key)
	})

	var sb strings.Builder
	for i := range kvs {
		sb.WriteString(kvs[i].value)
	}

	hash := sha256.New().Sum([]byte(sb.String()))

	hashEncr, err := base64.StdEncoding.DecodeString(tokenString)
	if err != nil {
		return false, err
	}

	data, err := s.privateKey.Decrypt(nil, hashEncr, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		return false, err
	}

	// s.logger.Info(context.Background(), "hashes", "in", string(hash), "out", string(data))

	return slices.Equal(hash, data), nil
}
