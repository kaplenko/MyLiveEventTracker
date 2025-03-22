package http

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// loadTemplates загружает HTML-шаблоны из директории templates.
func (h *Handler) loadTemplates() (*template.Template, error) {
	templates := template.New("")

	// Базовый путь к шаблонам (относительно расположения исполняемого файла)
	basePath := filepath.Join("pkg", "templates")

	// Шаблоны, которые нужно загрузить
	templateFiles := []string{
		"header.html",
		"footer.html",
		"login.html",
		"register.html",
		"home.html",
	}

	// Парсим каждый шаблон
	for _, file := range templateFiles {
		fullPath := filepath.Join(basePath, file)
		_, err := templates.ParseFiles(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", fullPath, err)
		}
	}

	return templates, nil
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// Загружаем шаблоны
	templates, err := h.loadTemplates()
	if err != nil {
		h.log.Error("Failed to load templates", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Данные для передачи в шаблон
	data := struct {
		Title string
		Error string
	}{
		Title: "Login",
		Error: "", // По умолчанию ошибки нет
	}

	// Рендерим шаблон
	err = templates.ExecuteTemplate(w, "login.html", data)
	if err != nil {
		h.log.Error("Failed to render template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	// Загружаем шаблоны
	templates, err := h.loadTemplates()
	if err != nil {
		h.log.Error("Failed to load templates", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Данные для передачи в шаблон
	data := struct {
		Title   string
		Message string
		Error   string
	}{
		Title:   "Register",
		Message: "", // По умолчанию сообщения нет
		Error:   "", // По умолчанию ошибки нет
	}

	// Рендерим шаблон
	err = templates.ExecuteTemplate(w, "register.html", data)
	if err != nil {
		h.log.Error("Failed to render template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	// Проверяем, авторизован ли пользователь
	userIDRaw := r.Context().Value("user_id")
	userID, ok := userIDRaw.(int64)
	if !ok {
		// Если пользователь не авторизован, перенаправляем на страницу входа
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Получаем данные пользователя из базы данных
	user, err := h.useCase.GetUser(r.Context(), userID)
	if err != nil {
		h.log.Error("Failed to get user", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Загружаем шаблоны
	templates, err := h.loadTemplates()
	if err != nil {
		h.log.Error("Failed to load templates", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Данные для передачи в шаблон
	data := struct {
		Title string
		User  struct {
			Name  string
			Email string
		}
	}{
		Title: "Home",
		User: struct {
			Name  string
			Email string
		}{
			Name:  user.Username,
			Email: user.Email,
		},
	}

	// Рендерим шаблон
	err = templates.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		h.log.Error("Failed to render template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
