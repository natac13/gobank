package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type APIServer struct {
	listAddr string
	store    Storage
}

func NewAPIServer(listAddr string, store Storage) *APIServer {
	return &APIServer{
		listAddr: listAddr,
		store:    store,
	}
}

/*
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	http.ListenAndServe(":3000", r)
}
*/

func (s *APIServer) Run() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/account/", makeHTTPHandleFunc(s.handleGetAccount))
	router.Get("/account", makeHTTPHandleFunc(s.handleGetAccount))
	router.Post("/account", makeHTTPHandleFunc(s.handleCreateAccount))

	router.Get("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID))
	router.Delete("/account/{id}", makeHTTPHandleFunc(s.handleDeleteAccount))
	router.Put("/transfer/{accountNumber}", makeHTTPHandleFunc(s.handleTransfer))

	log.Println("JSON API serverr is listening on: ", s.listAddr)
	http.ListenAndServe(s.listAddr, router)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	// db.get(id)
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

	account, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest) // &CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if _, err := s.store.CreateAccount(account); err != nil {
		return err
	}

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
	return WriteJSON(w, http.StatusOK, map[string]int{"deletedId": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}

	defer r.Body.Close()

	// fromAccountID, err := getID(r)
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type APIError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle the error
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return id, fmt.Errorf("invalid account id: %s", idString)
	}

	return id, nil
}
