package util

import (
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

// ExecuteRequest execute the request
func ExecuteRequest(router *mux.Router, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}
