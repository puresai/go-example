package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"demo7-docker-compose/endpoint"
)

var (
	ErrorBadReq = errors.New("params invalid")
)

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func MakeHttpHand(ctx context.Context, end *endpoint.UserEndpoints) http.Handler {
	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}
	r.Methods("POST").Path("/register").Handler(kithttp.NewServer(
		end.RegisterEndpoint,
		decodeRegisterRequest,
		encodeJSONResponse,
		options...,
	))

	r.Methods("POST").Path("/login").Handler(kithttp.NewServer(
		end.LoginEndpoint,
		decodeLoginRequest,
		encodeJSONResponse,
		options...,
	))

	return r
}
func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	if username == "" || password == "" || email == "" {
		return nil, ErrorBadReq
	}
	return &endpoint.RegisterRequest{
		Username: username,
		Password: password,
		Email:    email,
	}, nil
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		return nil, ErrorBadReq
	}
	return &endpoint.LoginRequest{
		Email:    email,
		Password: password,
	}, nil
}

func encodeJSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
