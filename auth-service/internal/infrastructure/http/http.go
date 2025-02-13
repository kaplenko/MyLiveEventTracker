package http

import (
	"auth-service/internal/entity"
	"auth-service/internal/usecase"
	"auth-service/pkg/jwt"
	"encoding/json"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

type Handler struct {
	usecase *usecase.UseCase
	router  *mux.Router
	log     *slog.Logger
}

type userDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func New(usecase *usecase.UseCase, r *mux.Router, log *slog.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		router:  r,
		log:     log,
	}
}

func (h *Handler) Router() *mux.Router {
	return h.router
}

func (h *Handler) SetupRoutes() {
	h.router.HandleFunc("/register", h.Register).Methods("POST")
	h.router.HandleFunc("/login", h.Login).Methods("POST")
	h.router.Handle("/profile", jwt.JWTMiddleware(http.HandlerFunc(h.Profile))).Methods("GET")
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req userDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	user := entity.User{
		ID:       req.ID,
		Username: req.Username,
		Email:    req.Email,
		PassHash: []byte(req.Password),
	}

	id, err := h.usecase.Registre(r.Context(), user)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]int64{"id": id}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error(err.Error())
		return
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req *userDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	token, err := h.usecase.Login(r.Context(), req.Email, []byte(req.Password))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.log.Info("User logged in")

	response := map[string]string{"token": token}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error(err.Error())
		return
	}
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("user_id")
	userID, ok := userIDRaw.(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.usecase.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		h.log.Error(err.Error())
		return
	}
}
