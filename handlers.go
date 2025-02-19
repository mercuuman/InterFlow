package main

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// Страница регистрации login GET
func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("login.html") // путь к странице регистрации
	if err != nil {
		log.Printf("Error parsing login template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Обработчик для регистрации пользователя signup POST
func SignUpPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json") // Устанавливаем заголовок Content-Type
	var req UserIn
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, `{"error": "Error hashing password"}`, http.StatusInternalServerError)
		return
	}

	req.Password = string(hashedPassword)
	token, err := RegisterUser(req)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrEmailExists) || errors.Is(err, ErrUsernameExists) {
			status = http.StatusConflict
		}

		w.WriteHeader(status)
		if jsonErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); jsonErr != nil {
			log.Printf("Error encoding JSON response: %v", jsonErr)
		}
		return
	}
	log.Printf("User %s successfully registered", req.Username)
	sendMail(req.Username, token)
	response := map[string]string{
		"message": "User registered successfully. Please check your email to complete the process.",
	}

	if jsonErr := json.NewEncoder(w).Encode(response); jsonErr != nil {
		log.Printf("Error encoding success response: %v", jsonErr)
	}
}

// !!! На доработке
func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}
	userID, err := VerifyEmail(token)
	if err != nil {
		status := http.StatusInternalServerError
		msg := "Database error"
		if errors.Is(err, ErrInvalidToken) {
			status = http.StatusUnauthorized
			msg = "Invalid or expired token"
		}
		http.Error(w, msg, status)
		log.Printf("Error verifying email: %v", err)
		return
	}

	accessToken, refreshToken, err := GenerateAccessRefresh(userID)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	// Сохраняем refresh-токен в cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		//Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	// Отправляем access-токен клиенту
	json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
}

// !!! на доработке
func RefreshAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем refresh-токен из cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Missing refresh token", http.StatusUnauthorized)
		return
	}

	// Обновляем только `access`-токен
	accessToken, err := RefreshAccessToken(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Отправляем новый `access`-токен клиенту
	json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
}
func RefreshTokensHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем refresh-токен из cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Missing refresh token", http.StatusUnauthorized)
		return
	}

	// Обновляем токены
	accessToken, refreshToken, err := RefreshTokens(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Устанавливаем новый refresh-токен, если он изменился
	if refreshToken != cookie.Value {
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			Expires:  time.Now().Add(7 * 24 * time.Hour),
		})
	}

	// Отправляем новый access-токен клиенту
	json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
}
