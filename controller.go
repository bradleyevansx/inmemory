package main

import (
	"encoding/json"
	"fmt"
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
	c.server.RegisterRoute("GET /{id}", makeHTTPFunc(c.handleGetById))
	c.server.RegisterRoute("POST /", makeHTTPFunc(c.handleInsert))
	c.server.RegisterRoute("DELETE /{id}", makeHTTPFunc(c.handleDelete))
}

func makeHTTPFunc(f apiServerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, err.Error())
		}
	}
}

func (s *APIController[T]) handleGet(w http.ResponseWriter, _ *http.Request) error {
	e, err := s.repository.Get() 
	if err != nil {
		return err
	}
	WriteJSON(w, http.StatusOK, e)
	return nil
}

func (s *APIController[T]) handleGetById(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	e, err := s.repository.GetById(id) 
	if err != nil {
		return err
	}
	WriteJSON(w, http.StatusOK, e)
	return nil
}

func (s *APIController[T]) handleInsert(w http.ResponseWriter, r *http.Request) error {
	var e T
	
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return fmt.Errorf("handleInsert: error decoding request body - %w", err)
	}
	defer r.Body.Close()
	 
	res, err := s.repository.Create(&e)
	if err != nil {
		return err
	}
	
	WriteJSON(w, http.StatusOK, res)

	return nil
}

func (s *APIController[T]) handleUpdate(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIController[T]) handleDelete(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	err := s.repository.Delete(id)
	if err != nil {
		return err
	}

	res := map[string]string{"message": "Entity deleted successfully"}

	WriteJSON(w, http.StatusOK, res)
	return nil
}