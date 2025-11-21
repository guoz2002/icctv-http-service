package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type baseResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func respondData(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(baseResponse{
		Success: true,
		Data:    data,
	})
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(baseResponse{
		Success: false,
		Error:   message,
	})
}

func respondResult(w http.ResponseWriter, err error, successStatus int, payload interface{}) {
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondData(w, successStatus, payload)
}

func decodeJSON(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return errors.New("empty body")
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return errors.New("empty body")
	}
	return json.Unmarshal(body, dst)
}

func parseIDFromQuery(r *http.Request) (int64, error) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return 0, errors.New("id is required")
	}
	return strconv.ParseInt(idStr, 10, 64)
}
