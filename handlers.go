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
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

//В логин пост хэндлере нужна проверка на isVerified

// Обработчик для регистрации пользователя signup POST
func SignUpPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleError(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	var req UserIn
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, errors.New("Invalid Request body"), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		handleError(w, errors.New("error hashing password"), http.StatusInternalServerError)
		return
	}

	req.Password = string(hashedPassword)
	token, err := RegisterUser(req)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrEmailExists) || errors.Is(err, ErrUsernameExists) {
			status = http.StatusConflict
		}
		handleError(w, err, status)
		return
	}

	log.Printf("User %s successfully registered", req.Username)
	sendMail(req.Username, token)

	response := map[string]string{
		"message": "User registered successfully. Please check your email to complete the process.",
	}
	json.NewEncoder(w).Encode(response)
}

// Отправка двух токенов при подтверждении почты
func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		handleError(w, errors.New("Missing token"), http.StatusBadRequest)
		return
	}
	userID, err := VerifyEmail(token)
	if err != nil {

		if errors.Is(err, ErrInvalidToken) {
			handleError(w, err, http.StatusUnauthorized)
		} else {
			handleError(w, err, http.StatusInternalServerError)
		}
		return
	}

	// Генерируем токены
	accessToken, err := GenerateJWT(userID, 15*time.Minute)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	refreshToken, err := GenerateJWT(userID, 7*24*time.Hour)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	// Сохраняем refresh-токен в cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		//Secure:   true,
		//SameSite: http.SameSiteStrictMode,
		Expires: time.Now().Add(7 * 24 * time.Hour),
	})

	// Отправляем access-токен клиенту
	json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
}

// !!! на доработке
func RefreshAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		handleError(w, errors.New("missing refresh token"), http.StatusUnauthorized)
		return
	}

	accessToken, err := RefreshAccessToken(cookie.Value)
	if err != nil {
		handleError(w, errors.New("invalid refresh token"), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
}

func RefreshTokensHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		handleError(w, errors.New("missing refresh token"), http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := RefreshTokens(cookie.Value)
	if err != nil {
		handleError(w, errors.New("invalid refresh token"), http.StatusUnauthorized)
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

	json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
}
