package transport

import "github.com/smtdfc/nagare/core/transport/http"

type Gateway struct {
	httpServer *http.Server
}

func (g *Gateway) Start() {
	go g.httpServer.Listen()
}

// @Injectable
func NewGateway(httpServer *http.Server) *Gateway {
	return &Gateway{httpServer: httpServer}
}
