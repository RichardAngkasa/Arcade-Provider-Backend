package middleware

import (
	"errors"
	"net/http"
)

type ctxKey int

const (
	CtxAdminID ctxKey = iota
	CtxClientID
	CtxRole
)

func MustClientID(r *http.Request) (int, error) {
	id, ok := r.Context().Value(CtxClientID).(int)
	if !ok {
		return 0, errors.New("unauthorized")
	}
	return id, nil
}

func MustAdminID(r *http.Request) (int, error) {
	id, ok := r.Context().Value(CtxAdminID).(int)
	if !ok {
		return 0, errors.New("unauthorized")
	}
	return id, nil
}
