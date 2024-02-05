package transport

import (
	"github.com/matheus-ds/day-trading-app/backend/internal/service"
)

type Transport interface {
}

func NewHTTPTransport(srv service.Service) *HTTPEndpoint {
	return &HTTPEndpoint{
		srv: srv,
	}
}

type HTTPEndpoint struct {
	srv service.Service
}
