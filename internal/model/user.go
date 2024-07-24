package model

import (
	"bind9-manager-service/internal/types"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func UserIsExist(db *sql.DB, username string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CreateUser(db *sql.DB, user types.CreateUserReq) error {
	isExist, err := UserIsExist(db, user.Username)
	if err != nil {
		return err
	}

	if isExist {
		return fmt.Errorf("user %s already exists", user.Username)
	} else {
		bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		if err != nil {
			return err
		}
		password := string(bytes)

		CreatedAt := time.Now().Format("2006-01-02 15:04:05")
		UpdatedAt := time.Now().Format("2006-01-02 15:04:05")
		_, err = db.Exec("INSERT INTO users (username, password, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", user.Username, password, user.Role, CreatedAt, UpdatedAt)
		return err
	}
}

func UpdateUserRole(db *sql.DB, user types.UpdateUserReq) error {
	isExist, err := UserIsExist(db, user.Username)
	if err != nil {
		return err
	}
	if isExist {
		UpdatedAt := time.Now().Format("2006-01-02 15:04:05")
		_, err = db.Exec("UPDATE users SET role = ?, updated_at = ? WHERE username = ?", user.Role, UpdatedAt, user.Username)
		return err
	} else {
		return fmt.Errorf("user %s does not exist", user.Username)
	}

}

func UpdateUserPass(db *sql.DB, user types.UpdateUserPassReq) error {
	isExist, err := UserIsExist(db, user.Username)
	if err != nil {
		return err
	}
	if isExist {
		bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		if err != nil {
			return err
		}
		password := string(bytes)

		UpdatedAt := time.Now().Format("2006-01-02 15:04:05")
		_, err = db.Exec("UPDATE users SET password = ?, updated_at = ? WHERE username = ?", password, UpdatedAt, user.Username)
		return err
	} else {
		return fmt.Errorf("user %s does not exist", user.Username)
	}
}

func DeleteUser(db *sql.DB, username string) error {
	isExist, err := UserIsExist(db, username)
	if err != nil {
		return err
	}

	if isExist {
		_, err := db.Exec("DELETE FROM users WHERE username = ?", username)
		return err
	} else {
		return fmt.Errorf("user %s does not exist", username)
	}
}

func GetUserByName(db *sql.DB, username string) (user types.GetUserResp, err error) {
	isExist, err := UserIsExist(db, username)
	if err != nil {
		return user, err
	}

	if isExist {
		err = db.QueryRow("SELECT id, username, role, created_at, updated_at FROM users WHERE username = ?", username).Scan(&user.Id, &user.Username, &user.Role, &user.CreateAt, &user.UpdateAt)
		return user, err
	} else {
		return user, fmt.Errorf("user %s does not exist", username)
	}
}

func GetPasswordByName(db *sql.DB, username string) (password string, err error) {
	isExist, err := UserIsExist(db, username)
	if err != nil {
		return password, err
	}
	if isExist {
		err = db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&password)
		return password, err
	} else {
		return password, fmt.Errorf("user %s does not exist", username)
	}
}

func GetAllUsers(db *sql.DB) (users []types.GetUserResp, err error) {
	rows, err := db.Query("SELECT id, username, role, created_at, updated_at FROM users")
	if err != nil {
		return users, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var user types.GetUserResp
		err = rows.Scan(&user.Id, &user.Username, &user.Role, &user.CreateAt, &user.UpdateAt)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return users, fmt.Errorf("no users found")
	}

	return users, nil
}
