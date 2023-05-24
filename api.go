package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	jwt "github.com/golang-jwt/jwt/v5"
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

	router.Get("/account/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleGetAccountByID), s.store))
	router.Delete("/account/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleDeleteAccount), s.store))
	router.Post("/transfer", makeHTTPHandleFunc(s.handleTransfer))

	router.Post("/login/{id}", makeHTTPHandleFunc(s.handleLogin))

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

// Creates a new account and returns it
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest) // &CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if _, err := s.store.CreateAccount(account); err != nil {
		return err
	}

	tokenStr, err := createJWTToken(account.ID)
	if err != nil {
		return err
	}

	fmt.Printf("tokenStr: %s\n", tokenStr)

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
	return WriteJSON(w, http.StatusOK, transferReq)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	// get the account id from the urlAccountId
	accountID, err := getID(r)

	if err != nil {
		return err
	}

	// find the account in the db
	account, err := s.store.GetAccountByID(accountID)

	if err != nil {
		return err
	}

	if account == nil || account.ID != accountID {
		return fmt.Errorf("account not found")
	}

	jwtToken, err := createJWTToken(account.ID)

	if err != nil {
		return err
	}

	// if found then create and JWT token and return it to the user
	return WriteJSON(w, http.StatusOK, map[string]string{"token": jwtToken})
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

func unauthorizedError(w http.ResponseWriter) {
	WriteJSON(w, http.StatusUnauthorized, APIError{Error: "Unauthorized"})
}

// This code creates a handler function that checks if the request has a valid JWT token.
// It returns an error if the token is not valid.
// If the token is valid, it calls the handlerFunc.
func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check if the request has a valid JWT token
		fmt.Println("withJWTAuth: func")
		tokenStr := r.Header.Get("x-jwt-token")
		token, err := validateJWTToken(tokenStr)

		if err != nil || !token.Valid {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		claimsAccountID := claims["account_id"].(float64)
		urlAccountId, err := getID(r)
		if err != nil {
			// use generic error message for security reasons
			unauthorizedError(w)
			return
		}

		if int(claimsAccountID) != urlAccountId {
			unauthorizedError(w)
			return
		}

		account, err := s.GetAccountByID(int(claimsAccountID))

		if err != nil {
			unauthorizedError(w)
			return
		}

		if account == nil || account.ID != int(claimsAccountID) {
			unauthorizedError(w)
			// log to central logging what is really going on in each of the errors
			// ie send to AWS CloudWatch , sentry, etc
			return
		}

		fmt.Printf("claims: %+v\n", claims)
		fmt.Printf("claimsID: %+v\n", claimsAccountID)

		// if not return an error
		// if valid call the handlerFunc
		handlerFunc(w, r)
	}
}

// validateJWTToken checks if the JWT token is valid
func validateJWTToken(tokenStr string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET env variable not set")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unauthorized")
			// return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Unauthorized")
	}

	return token, nil

}

// createJWTToken creates a JWT token for the given account ID.
// The token is signed using the JWT_SECRET env variable.
func createJWTToken(id int) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET env variable not set")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account_id": id,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(secret))
}
