package main

import (
	"log"
	"net/http"
)

type Server struct{
	mux *http.ServeMux
	listenAddr string
}

func NewServer(listenAddr string) *Server{
	mux := http.NewServeMux()
	return &Server{
		mux: mux, listenAddr: listenAddr,
	}
}

func (s *Server)Run(){
	http.ListenAndServe(s.listenAddr, s.mux)
    log.Println("Server started. Listening on port :3000")
}

func (s *Server)RegisterRoute(route string, f http.HandlerFunc){
	s.mux.HandleFunc(route, f)
}