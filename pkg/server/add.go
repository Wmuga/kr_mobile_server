package server

import (
	"encoding/json"
	"io"
	"net/http"
)

func (s *Server) add(w http.ResponseWriter, r *http.Request) {
	in, err := addDecodeRequest(r)
	if err != nil {
		transportError(w, http.StatusBadRequest, err)
		return
	}
	resp, err := s.service.AddPass(r.Context(), in.Mac, in.Time)

	if err != nil {
		transportError(w, http.StatusInternalServerError, err)
		return
	}

	transportResponse(w, resp, 1)
}

func addDecodeRequest(r *http.Request) (res addInput, err error) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return res, err
	}
	defer r.Body.Close()
	err = json.Unmarshal(bytes, &res)
	return res, err
}
