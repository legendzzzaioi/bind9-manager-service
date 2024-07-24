package model

import (
	"database/sql"
	"time"

	"bind9-manager-service/internal/types"
)

func CreateOperationLog(db *sql.DB, username string, operation string, context string) error {

	_, err := db.Exec("INSERT INTO operation_logs (username, operation, context, created_at) VALUES (?, ?, ?, ?)", username, operation, context, time.Now().Format("2006-01-02 15:04:05"))

	return err
}

func CreateUserLog(db *sql.DB, username string, operation string, context string) error {
	_, err := db.Exec("INSERT INTO user_logs (username, operation, context, created_at) VALUES (?, ?, ?, ?)", username, operation, context, time.Now().Format("2006-01-02 15:04:05"))

	return err
}

func CreateLoginLog(db *sql.DB, username string, ip string, operation string) error {

	_, err := db.Exec("INSERT INTO login_logs (username, ip, operation, created_at) VALUES (?, ?, ?, ?)", username, ip, operation, time.Now().Format("2006-01-02 15:04:05"))

	return err

}

// GetOperationLog retrieves operation logs with pagination support.
func GetOperationLog(db *sql.DB, page, pageSize int) ([]types.OperationLog, error) {
	offset := (page - 1) * pageSize
	query := "SELECT id, username, operation, context, created_at FROM operation_logs ORDER BY created_at DESC LIMIT ? OFFSET ?"
	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var logs []types.OperationLog
	for rows.Next() {
		var log types.OperationLog
		err := rows.Scan(&log.Id, &log.Username, &log.Operation, &log.Context, &log.CreateAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetOperationLogCount retrieves the total count of operation logs.
func GetOperationLogCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM operation_logs").Scan(&count)
	return count, err
}

// GetUserLog retrieves user logs with pagination support.
func GetUserLog(db *sql.DB, page, pageSize int) ([]types.UserLog, error) {
	offset := (page - 1) * pageSize
	query := "SELECT id, username, operation, context, created_at FROM user_logs ORDER BY created_at DESC LIMIT ? OFFSET ?"
	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var logs []types.UserLog
	for rows.Next() {
		var log types.UserLog
		err := rows.Scan(&log.Id, &log.Username, &log.Operation, &log.Context, &log.CreateAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetUserLogCount retrieves the total count of user logs.
func GetUserLogCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user_logs").Scan(&count)
	return count, err
}

// GetLoginLog retrieves login logs with pagination support.
func GetLoginLog(db *sql.DB, page, pageSize int) ([]types.LoginLog, error) {
	offset := (page - 1) * pageSize
	query := "SELECT id, username, ip, operation, created_at FROM login_logs ORDER BY created_at DESC LIMIT ? OFFSET ?"
	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var logs []types.LoginLog
	for rows.Next() {
		var log types.LoginLog
		err := rows.Scan(&log.Id, &log.Username, &log.Ip, &log.Operation, &log.CreateAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetLoginLogCount retrieves the total count of login logs.
func GetLoginLogCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM login_logs").Scan(&count)
	return count, err
}
