package http

import (
	"auth-service/internal/entity"
	"auth-service/internal/usecase"
	"encoding/json"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

type Handler struct {
	usecase usecase.UseCase
	lg      *slog.Logger
}

func NewHandler(usecase usecase.UseCase, lg *slog.Logger) *Handler {
	return &Handler{usecase, lg}
}

func (h *Handler) SetupRoutes() {
	r := mux.NewRouter()

	r.HandleFunc("/register", h.Register).Methods("POST")
	r.HandleFunc("login", h.Login).Methods("POST")
	r.HandleFunc("/profile", h.Profile).Methods("GET")
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req entity.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	id, err := h.usecase.Registre(r.Context(), req)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	response := map[string]int64{"id": id}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req entity.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	token, err := h.usecase.Login(r.Context(), req.Email, req.PasswordHash)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	h.lg.Info("User logged in", "email", req.Email)

	response := map[string]string{"token": token}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)

	user, err := h.usecase.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
