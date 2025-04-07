package storage

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
)

type User struct {
	ID       int
	Username string
	Password string
	Email    string
}

func UserExists(username string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username=?)"

	err := DB.QueryRow(query, username).Scan(&exists)

	if err != nil {
		return false, err
	}
	return exists, nil
}

func CreateUser(username, password, email string) error {
	query := "INSERT INTO users (username, password, email) VALUES (?, ?, ?)"
	_, err := DB.Exec(query, username, password, email)

	if err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}
	return nil
}

func ValidateUser(username, password string) (bool, error) {
	var dbPassword string
	query := "SELECT password FROM users WHERE username=?"
	err := DB.QueryRow(query, username).Scan(&dbPassword)

	if err == sql.ErrNoRows {
		log.Println("❌ Пользователь не найден:", username)
		return false, nil
	} else if err != nil {
		log.Println("❌ Ошибка при запросе пользователя:", err)
		return false, err
	}

	log.Println("✅ Найден пользователь:", username, "| Пароль в БД:", dbPassword)
	log.Println("Введённый пароль:", password)

	return dbPassword == password, nil
}
