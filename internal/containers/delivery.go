package containers

import "net/http"

type Handlers interface {
	GetAll() http.HandlerFunc
	SearchByIP() http.HandlerFunc
	SetAll() http.HandlerFunc
	GetHistory() http.HandlerFunc
}
