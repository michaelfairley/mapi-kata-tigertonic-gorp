package web

import (
	"net/http"
	"net/url"
)

type tigerTonicHandler func(*url.URL, http.Header, interface{}) (int, http.Header, interface{}, error)
