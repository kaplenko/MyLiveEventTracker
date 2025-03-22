package http

import (
	"auth-service/internal/entity"
	"auth-service/internal/usecase"
	"auth-service/pkg/jwt"
	"auth-service/pkg/oauth2/github"
	"auth-service/pkg/oauth2/google"
	"encoding/json"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

type Handler struct {
	useCase       *usecase.UseCase
	githubService *github.Service
	googleService *google.Service
	router        *mux.Router
	log           *slog.Logger
}

type userDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func New(useCase *usecase.UseCase, gh *github.Service, google *google.Service, r *mux.Router, log *slog.Logger) *Handler {
	return &Handler{
		useCase:       useCase,
		githubService: gh,
		googleService: google,
		router:        r,
		log:           log,
	}
}

func (h *Handler) Router() *mux.Router {
	return h.router
}

func (h *Handler) SetupRoutes() {
	h.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("pkg/static"))))

	h.router.HandleFunc("/login", h.LoginPage).Methods("GET")
	h.router.HandleFunc("/register", h.RegisterPage).Methods("GET")
	h.router.Handle("/", jwt.JWTMiddleware(http.HandlerFunc(h.HomePage))).Methods("GET")

	h.router.Handle("/profile", jwt.JWTMiddleware(http.HandlerFunc(h.Profile))).Methods("GET")

	h.router.HandleFunc("/register", h.Register).Methods("POST")
	h.router.HandleFunc("/login", h.Login).Methods("POST")

	h.router.HandleFunc("/auth/github", h.GithubLoginRedirect).Methods("GET")
	h.router.HandleFunc("/auth/github/callback", h.GithubCallback).Methods("GET")

	h.router.HandleFunc("/auth/google", h.GoogleLoginRedirect).Methods("GET")
	h.router.HandleFunc("/auth/google/callback", h.GoogleCallback).Methods("GET")
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	// Парсим данные из формы
	name := r.FormValue("_name")
	email := r.FormValue("_email")
	password := r.FormValue("_password")

	// Создаем объект пользователя
	user := entity.User{
		Username: name,
		Email:    email,
		PassHash: []byte(password), // Пароль нужно хэшировать перед сохранением в БД
	}

	// Регистрируем пользователя через useCase
	_, err := h.useCase.Registre(r.Context(), user)
	if err != nil {
		// Если ошибка, показываем страницу регистрации с сообщением об ошибке
		data := struct {
			Title   string
			Message string
			Error   string
		}{
			Title:   "Register",
			Message: "",
			Error:   "Registration failed: " + err.Error(),
		}

		templates, err := h.loadTemplates()
		if err != nil {
			h.log.Error("Failed to load templates", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = templates.ExecuteTemplate(w, "register.html", data)
		if err != nil {
			h.log.Error("Failed to render template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// Если регистрация успешна, показываем сообщение об успехе
	data := struct {
		Title   string
		Message string
		Error   string
	}{
		Title:   "Register",
		Message: "Registration successful! Please login.",
		Error:   "",
	}

	templates, err := h.loadTemplates()
	if err != nil {
		h.log.Error("Failed to load templates", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "register.html", data)
	if err != nil {
		h.log.Error("Failed to render template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("_email")
	password := r.FormValue("_password")

	token, err := h.useCase.Login(r.Context(), email, []byte(password))
	if err != nil {
		data := struct {
			Title string
			Error string
		}{
			Title: "Login",
			Error: "Invalid email or password",
		}

		templates, err := h.loadTemplates()
		if err != nil {
			h.log.Error("Failed to load templates", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = templates.ExecuteTemplate(w, "login.html", data)
		if err != nil {
			h.log.Error("Failed to render template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	})
	http.Redirect(w, r, "/home", http.StatusFound)
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("user_id")
	userID, ok := userIDRaw.(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.useCase.GetUser(r.Context(), userID)
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
