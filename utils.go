package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	jwtSecret = []byte("supersecretkey")
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateTokens создает accessToken и refreshToken
func GenerateAccessRefresh(userID int) (string, string, error) {
	accessToken, err := createJWTToken(userID, 15*time.Minute) // Access токен на 15 минут
	if err != nil {
		return "", "", err
	}
	refreshToken, err := createJWTToken(userID, 7*24*time.Hour) // Refresh токен на 7 дней
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// !!!доработка
func createJWTToken(userID int, duration time.Duration) (string, error) {
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

// !!!доработка
func RefreshAccessToken(refreshToken string) (string, error) {
	// Разбираем refresh-токен
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	// Получаем userID из токена
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID := claims.UserID

	// Генерируем новый access-токен
	accessToken, err := createJWTToken(userID, 15*time.Minute)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

// !!!доработка
func RefreshTokens(refreshToken string) (string, string, error) {
	// Разбираем refresh-токен
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	// Получаем userID из токена
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}
	userID := claims.UserID

	// Проверяем, истекает ли refresh-токен (например, если осталось <24 часов)
	timeLeft := claims.ExpiresAt.Time.Sub(time.Now())
	if timeLeft > 24*time.Hour {
		// Если до истечения refresh-токена больше 24 часов — возвращаем новый access-токен, но оставляем старый refresh-токен
		accessToken, err := createJWTToken(userID, 15*time.Minute)
		if err != nil {
			return "", "", err
		}
		return accessToken, refreshToken, nil
	}

	// Генерируем новые токены
	accessToken, err := createJWTToken(userID, 15*time.Minute)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := createJWTToken(userID, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}
