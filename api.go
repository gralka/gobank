package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
  listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
  return &APIServer{
    listenAddr: listenAddr,
  }
}

func (s *APIServer) Run() {
  router := mux.NewRouter()

  router.HandleFunc("/accounts/{id}", makeHTTPHandleFunc(s.handleAccount));
  router.HandleFunc("/accounts", makeHTTPHandleFunc(s.handleAccount));

  log.Println("JSON API Service Running on port: ", s.listenAddr)

  http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
  switch r.Method {
  case "GET":
    return s.handleGetAccount(w, r)
  case "POST":
    return s.handleCreateAccount(w, r)
  case "DELETE":
    return s.handleDeleteAccount(w, r)
  default:
    return fmt.Errorf("method not allowed %s", r.Method)
  }
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
  id := mux.Vars(r)["id"]

  // if _, ok := vars["id"]; ok {
    // return WriteJSON(w, http.StatusOK, vars)
  // }
  fmt.Println(id)

  // return WriteJSON(w, http.StatusOK, []string{"account1", "account2"})
  return WriteJSON(w, http.StatusOK, Account{})
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
  return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
  return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
  return nil
}

func WriteJSON(w http.ResponseWriter, status int, value any) error {
  w.Header().Add("Content-Type", "application/json")
  w.WriteHeader(status)

  return json.NewEncoder(w).Encode(value)
}

type apiFunc func (http.ResponseWriter, *http.Request) error

type ApiError struct {
  Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
    if err := f(w, r); err != nil {
      WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
    }
  }
}

