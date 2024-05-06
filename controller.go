package main

import (
	"encoding/json"
	"net/http"
)

type APIController[T IEntity] struct{
	server *Server
	repository IRepository[T]
}

type apiServerFunc func (w http.ResponseWriter, r *http.Request) error 

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(v); err != nil {
		return err
	}
	return nil
}

func NewAPIController[T IEntity](server *Server, repository IRepository[T])*APIController[T]{
	return &APIController[T]{server: server, repository: repository}
}

func (c *APIController[T]) RegisterRoutes(){
	c.server.RegisterRoute("GET /", makeHTTPFunc(c.handleGet))
}

func makeHTTPFunc(f apiServerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, err.Error())
		}
	}
}

func (s *APIController[T]) handleGet(w http.ResponseWriter, r *http.Request) error {
	e, err := s.repository.Get() 
	if err != nil {
		return err
	}
	WriteJSON(w, http.StatusOK, e)
	return nil
}

func (s *APIController[T]) handleInsert(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIController[T]) handleUpdate(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIController[T]) handleDelete(w http.ResponseWriter, r *http.Request) error {
	return nil
}