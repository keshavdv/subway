package context

import (
	"github.com/gorilla/context"
	"github.com/hashicorp/yamux"
	"github.com/unrolled/render"
	"net/http"
)

type SubwayMiddleware struct {
	Mux    *yamux.Session
	Render *render.Render
}

func CreateSubway(session *yamux.Session, render *render.Render) *SubwayMiddleware {
	return &SubwayMiddleware{session, render}
}

func (session SubwayMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Mux", session.Mux)
	context.Set(r, "Render", session.Render)
	next(rw, r)
}
