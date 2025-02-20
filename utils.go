package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
)

var (
	jwtSecret = []byte("supersecretkey")
)

func handleError(w http.ResponseWriter, err error, status int) {
	log.Println("Error:", err) // Логирование ошибки
	w.WriteHeader(status)
	response := map[string]string{"error": err.Error()}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// токен для uRL в письме
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func GenerateJWT(userID int, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Обновление access-токена по refresh-токену
func RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := ParseRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	return GenerateJWT(claims.UserID, 15*time.Minute)
}

// Обновление пары токенов
func RefreshTokens(refreshToken string) (string, string, error) {
	claims, err := ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	timeLeft := time.Until(claims.ExpiresAt.Time)
	if timeLeft > 24*time.Hour {
		accessToken, err := GenerateJWT(claims.UserID, 15*time.Minute)
		return accessToken, refreshToken, err
	}

	// Генерируем новые токены
	accessToken, err := GenerateJWT(claims.UserID, 15*time.Minute)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := GenerateJWT(claims.UserID, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

// Функция для разбора refresh-токена
func ParseRefreshToken(refreshToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
