package utils

import (
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"sofa-logs-servers/infra/zincsearch"
)

type ErrorForm struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func Middleware(
	next func(w http.ResponseWriter, r *http.Request, zincClient zincsearch.ZincClient), zincClient zincsearch.ZincClient,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, zincClient)
	}
}

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func WriteErr(w http.ResponseWriter, err string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	PanicErr(jsoniter.NewEncoder(w).Encode(ErrorForm{Message: err, StatusCode: statusCode}))
}

func WriteJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	PanicErr(jsoniter.NewEncoder(w).Encode(data))
}
