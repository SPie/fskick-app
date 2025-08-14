package server

import (
	"embed"
	"fmt"
	"net/http"
)

type Server struct {
	mux  *http.ServeMux
	addr string
}

func New(addr string) *Server {
	return &Server{
		mux:  http.NewServeMux(),
		addr: addr,
	}
}

func (server *Server) Get(route string, handler func(http.ResponseWriter, *http.Request)) {
	server.mux.HandleFunc(fmt.Sprintf("GET %s", route), handler)
}

func (server *Server) Post(route string, handler func(http.ResponseWriter, *http.Request)) {
	server.mux.HandleFunc(fmt.Sprintf("POST %s", route), handler)
}

func (server *Server) HandleStatic(static embed.FS) {
	server.mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(static)))
}

func (server *Server) Run() error {
	return http.ListenAndServe(server.addr, server.mux)
}
