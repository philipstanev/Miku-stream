package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/philipstanev/Miku-stream/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) AuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /login/", s.loginHandler)
	mux.HandleFunc("POST /register/", s.registerHandler)
	mux.HandleFunc("GET /secretMessage/", s.secretMessage)
	mux.HandleFunc("GET /proxyAuth", s.proxyAuthRequest)
}

func (s *Service) proxyAuthRequest(w http.ResponseWriter, r *http.Request) {
	err := middleware.CheckPermission(w, r, s.S)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *Service) secretMessage(w http.ResponseWriter, r *http.Request) {
	err := middleware.CheckPermission(w, r, s.S)
	if err != nil {
		return
	}

	res := secretMessageResponse{Message: "Dyuzov"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

type secretMessageResponse struct {
	Message string `json:"message"`
}

func (s *Service) registerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	userID, err := s.Register(ctx, req.Password, req.Username)
	if err != nil {
		http.Error(w, "Couldn't create user", http.StatusInternalServerError)
		return
	}
	res := registerResponse{UserID: userID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)

}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerResponse struct {
	UserID string `json:"userID"`
}

func (s *Service) loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	sessionID, err := s.Login(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			http.Error(w, "Wrong password", 401)
			return
		} else if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "No such user", 404)
			return
		}
	}
	res := loginResponse{SessionID: sessionID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(res)
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	SessionID string `json:"sessionID"`
}
