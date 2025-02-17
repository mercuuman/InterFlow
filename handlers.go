package main

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
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
