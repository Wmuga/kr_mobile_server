package server

import (
	"encoding/json"
	"io"
	"net/http"
)

func (s *Server) listAll(w http.ResponseWriter, r *http.Request) {
	in, err := listDecodeRequest(r)
	if err != nil {
		transportError(w, http.StatusBadRequest, err)
		return
	}

	resp, err := s.service.ListAll(r.Context(), in.Offset)
	if err != nil {
		transportError(w, http.StatusInternalServerError, err)
		return
	}

	transportResponse(w, resp, len(resp))
}

func (s *Server) listToday(w http.ResponseWriter, r *http.Request) {
	in, err := listDecodeRequest(r)
	if err != nil {
		transportError(w, http.StatusBadRequest, err)
		return
	}

	resp, err := s.service.ListToday(r.Context(), in.Offset)
	if err != nil {
		transportError(w, http.StatusInternalServerError, err)
		return
	}

	transportResponse(w, resp, len(resp))
}

func listDecodeRequest(r *http.Request) (res listInput, err error) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return res, err
	}
	defer r.Body.Close()
	err = json.Unmarshal(bytes, &res)
	return res, err
}
