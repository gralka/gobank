package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type APIServer struct {
  listenAddr string
  store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
  return &APIServer{
    listenAddr: listenAddr,
    store:      store,
  }
}

func (s *APIServer) Run() {
  router := mux.NewRouter()

  router.HandleFunc("/accounts/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleGetAccountByID), s.store));
  router.HandleFunc("/accounts", makeHTTPHandleFunc(s.handleAccount));
  router.HandleFunc("/transfers", makeHTTPHandleFunc(s.handleTransfer));

  log.Println("JSON API Service Running on port: ", s.listenAddr)

  http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
  switch r.Method {
  case "GET":
    return s.handleGetAllAccounts(w, r)
  case "POST":
    return s.handleCreateAccount(w, r)
  default:
    return fmt.Errorf("method not allowed %s", r.Method)
  }
}

func (s *APIServer) handleGetAllAccounts(w http.ResponseWriter, r *http.Request) error {
  accounts, err := s.store.GetAccounts()

  if err != nil {
    return err
  }

  return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
  id, err := getID(r)

  if err != nil {
    return err
  }

  switch r.Method {
  case "GET":
    account, err := s.store.GetAccountByID(id)

    if err != nil {
      return err
    }

    return WriteJSON(w, http.StatusOK, account)
  case "DELETE":
    return s.handleDeleteAccount(w, r)
  default:
    return fmt.Errorf("method not allowed %s", r.Method)
  }
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
  createAccountReq := new(CreateAccountRequest)

  if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
    return err
  }

  r.Body.Close()

  account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)

  newAccount, err := s.store.CreateAccount(account)

  if err != nil {
    return err
  }

  tokenString, err := createJWT(newAccount)

  if err != nil {
    return err
  }


  fmt.Println("tokenString: ", tokenString)

  return WriteJSON(w, http.StatusCreated, account) 
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
  id, err := getID(r)

  if err != nil {
    return err
  }

  if err := s.store.DeleteAccount(id); err != nil {
    return err
  }

  return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id}) 
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
  if r.Method != "POST" {
    return fmt.Errorf("method not allowed %s", r.Method)
  }

  transferReq := new(CreateTransferRequest)

  if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
    return err
  }

  r.Body.Close()

  return WriteJSON(w, http.StatusOK, transferReq)
}

func WriteJSON(w http.ResponseWriter, status int, value any) error {
  w.Header().Add("Content-Type", "application/json")
  w.WriteHeader(status)

  return json.NewEncoder(w).Encode(value)
}

func permissionDenied(w http.ResponseWriter) {
  WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "permission denied"})
}

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
    tokenString := r.Header.Get("x-jwt-token")

    token, err := validateJWT(tokenString)

    if err != nil || !token.Valid {
      permissionDenied(w)
      return
    }

    id, err := getID(r)

    if err != nil {
      permissionDenied(w)
      return
    }

    account, err := s.GetAccountByID(id)

    if err != nil {
      permissionDenied(w)
      return
    }

    claims := token.Claims.(jwt.MapClaims)

    if account.Number != int(claims["accountNumber"].(float64)) {
      permissionDenied(w)
      return
    }

    handlerFunc(w, r)
  }
}

func createJWT(account *Account) (string, error) {
  claims := &jwt.MapClaims{
    "accountNumber": account.Number,
  }

  secret := os.Getenv("JWT_SECRET")
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

  return token.SignedString([]byte(secret))
}
  
func validateJWT(tokenString string) (*jwt.Token, error) {
  secret := os.Getenv("JWT_SECRET")

  return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf( "Unexpected signing method: %v", token.Header["alg"])
    }

    return []byte(secret), nil
  })
}

type apiFunc func (http.ResponseWriter, *http.Request) error

type ApiError struct {
  Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
    if err := f(w, r); err != nil {
      WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
    }
  }
}

func getID(r *http.Request) (int, error) {
  idStr := mux.Vars(r)["id"]

  id, err := strconv.Atoi(idStr)

  if err != nil {
    return 0, fmt.Errorf("invalid id given: %s", idStr)
  }

  return id, nil
}
