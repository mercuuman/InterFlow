package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
)

// Страница регистрации login GET
func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("login.html") // путь к странице регистрации
	if err != nil {
		fmt.Println(err)
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
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
