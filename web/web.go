package web

import (
	"net/http"
	"net/url"
)

type tigerTonicHandler func(*url.URL, http.Header, interface{}) (int, http.Header, interface{}, error)

func validationError(errors map[string][]string) (int, http.Header, interface{}, error) {
	fullErrors := map[string]map[string][]string{
		"errors": errors,
	}

	return 422, http.Header{}, fullErrors, nil
}
