package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

var (
	ErrEmailExists    = errors.New("email already exists")
	ErrUsernameExists = errors.New("username already exists")
	db                *pgxpool.Pool
)

func initDB() error {
	connStr := "postgresql://postgres:641@localhost:5432/task_management"
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	db = pool
	return nil
}
func closeDB() {
	if db != nil {
		db.Close()
	}
}

// User registration блок
func insertUser(ctx context.Context, tx pgx.Tx, user UserIn) (int, error) {
	var userID int
	query := `INSERT INTO users (Username, Email, PasswordHash, CreatedAt) 
              VALUES ($1, $2, $3, $4) RETURNING id`
	err := tx.QueryRow(ctx, query, user.Username, user.Email, user.Password, time.Now()).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Код SQLSTATE для уникального ограничения
			switch pgErr.ConstraintName {
			case "users_email_key":
				return 0, ErrEmailExists
			case "users_username_key":
				return 0, ErrUsernameExists
			}
		}
		return 0, fmt.Errorf("Ошибка при добавлении пользователя: %v", err)
	}
	return userID, nil
}
func createVerificationToken(ctx context.Context, tx pgx.Tx, userID int) (string, error) {
	// Генерация токена
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	// Сохранение токена в базу данных
	query := `INSERT INTO user_state (UserID, token, isVerified, CreatedAt) 
              VALUES ($1, $2, FALSE, $3)`
	_, err = tx.Exec(ctx, query, userID, token, time.Now())
	if err != nil {
		return "", err
	}

	return token, nil
}
func RegisterUser(user UserIn) (string, error) {
	ctx := context.Background()

	// Создание транзакции с использованием pgx.Tx
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", fmt.Errorf("error starting transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	userId, err := insertUser(ctx, tx, user)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return "", err
	}
	token, err := createVerificationToken(ctx, tx, userId)
	if err != nil {
		log.Printf("Error creating verification token: %v", err)
		return "", err
	}
	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("error committing transaction: %w", err)
	}
	return token, nil
}
