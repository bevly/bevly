package httpfilestub

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func ServerValidated(file string, validator func(*http.Request) bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if validator != nil && !validator(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			w.Write(bytes)
		}
	}))
}

func Server(file string) *httptest.Server {
	return ServerValidated(file, nil)
}
