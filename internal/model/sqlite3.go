package model

import (
	"bind9-manager-service/internal/types"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
)

// record 类型
var validTypes = map[string]bool{
	"A":        true,
	"AAAA":     true,
	"CAA":      true,
	"CNAME":    true,
	"DNSKEY":   true,
	"IPSECKEY": true,
	"KEY":      true,
	"MX":       true,
	"NS":       true,
	"PTR":      true,
	"SPF":      true,
	"SRV":      true,
	"TLSA":     true,
	"TXT":      true,
}

// 通用事务处理函数
func withTransaction(db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

// ZoneIsExist 检查 zone 是否存在
func ZoneIsExist(tx *sql.Tx, domain string) (bool, error) {
	var exists bool
	if domain == "" {
		return exists, fmt.Errorf("domain cannot be empty")
	}
	query := "SELECT EXISTS(SELECT 1 FROM zones WHERE domain=?)"
	err := tx.QueryRow(query, domain).Scan(&exists)
	return exists, err
}

// RecordIsExist 检查 record 是否存在
func RecordIsExist(tx *sql.Tx, id int) (bool, error) {
	var exists bool
	if id <= 0 {
		return exists, fmt.Errorf("id must be a positive integer")
	}
	query := "SELECT EXISTS(SELECT 1 FROM records WHERE id=?)"
	err := tx.QueryRow(query, id).Scan(&exists)
	return exists, err
}

// GetZoneByDomain 获取指定 zone
func GetZoneByDomain(db *sql.DB, domain string) (types.Zone, error) {
	var zone types.Zone
	err := withTransaction(db, func(tx *sql.Tx) error {
		exists, err := ZoneIsExist(tx, domain)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("domain %s does not exist", domain)
		}

		query := "SELECT domain, soa_ttl, soa_cache_ttl, soa_expire, soa_mail_address, soa_primary_name_server, soa_refresh, soa_retry, soa_serial FROM zones WHERE domain = ?"
		row := tx.QueryRow(query, domain)
		err = row.Scan(&zone.Domain, &zone.Ttl, &zone.CacheTtl, &zone.Expire, &zone.MailAddress, &zone.PrimaryNameServer, &zone.Refresh, &zone.Retry, &zone.Serial)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return zone, err
	}
	return zone, nil
}

// GetZones 获取 zones 列表
func GetZones(db *sql.DB) ([]types.Zone, error) {
	rows, err := db.Query("SELECT domain, soa_ttl, soa_cache_ttl, soa_expire, soa_mail_address, soa_primary_name_server, soa_refresh, soa_retry, soa_serial FROM zones")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var zones []types.Zone
	for rows.Next() {
		var z types.Zone
		err := rows.Scan(&z.Domain, &z.Ttl, &z.CacheTtl, &z.Expire, &z.MailAddress, &z.PrimaryNameServer, &z.Refresh, &z.Retry, &z.Serial)
		if err != nil {
			return nil, err
		}
		zones = append(zones, z)
	}

	if len(zones) == 0 {
		return nil, fmt.Errorf("no zones found")
	}

	return zones, nil
}

// CreateZone 创建 zone
func CreateZone(db *sql.DB, zone types.ZoneReq) error {
	return withTransaction(db, func(tx *sql.Tx) error {
		exists, err := ZoneIsExist(tx, zone.Domain)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("domain %s already exists", zone.Domain)
		}

		stmt, err := tx.Prepare("INSERT INTO zones(domain, soa_ttl, soa_cache_ttl, soa_expire, soa_mail_address, soa_primary_name_server, soa_refresh, soa_retry, soa_serial) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer func(stmt *sql.Stmt) {
			_ = stmt.Close()
		}(stmt)

		// 获取当前时间戳
		serial := time.Now().Unix()
		_, err = stmt.Exec(zone.Domain, zone.Ttl, zone.CacheTtl, zone.Expire, zone.MailAddress, zone.PrimaryNameServer, zone.Refresh, zone.Retry, serial)
		return err
	})
}

// UpdateZone 更新 zone
func UpdateZone(db *sql.DB, zone types.ZoneReq) error {
	return withTransaction(db, func(tx *sql.Tx) error {
		exists, err := ZoneIsExist(tx, zone.Domain)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("domain %s does not exist", zone.Domain)
		}

		stmt, err := tx.Prepare("UPDATE zones SET soa_ttl=?, soa_cache_ttl=?, soa_expire=?, soa_mail_address=?, soa_primary_name_server=?, soa_refresh=?, soa_retry=?, soa_serial=? WHERE domain=?")
		if err != nil {
			return err
		}
		defer func(stmt *sql.Stmt) {
			_ = stmt.Close()
		}(stmt)

		// 获取当前时间戳
		serial := time.Now().Unix()
		_, err = stmt.Exec(zone.Ttl, zone.CacheTtl, zone.Expire, zone.MailAddress, zone.PrimaryNameServer, zone.Refresh, zone.Retry, serial, zone.Domain)
		return err
	})
}

// DeleteZone 删除 zone
func DeleteZone(db *sql.DB, domain string, record bool) error {
	return withTransaction(db, func(tx *sql.Tx) error {
		exists, err := ZoneIsExist(tx, domain)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("domain %s does not exist", domain)
		}

		stmt, err := tx.Prepare("DELETE FROM zones WHERE domain=?")
		if err != nil {
			return err
		}
		defer func(stmt *sql.Stmt) {
			_ = stmt.Close()
		}(stmt)
		_, err = stmt.Exec(domain)
		if err != nil {
			return err
		}

		if record {
			stmt, err := tx.Prepare("DELETE FROM records WHERE domain=?")
			if err != nil {
				return err
			}
			defer func(stmt *sql.Stmt) {
				_ = stmt.Close()
			}(stmt)
			_, err = stmt.Exec(domain)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// GetRecords 获取 records 列表
func GetRecords(db *sql.DB, domain string) ([]types.Record, error) {
	var records []types.Record

	err := withTransaction(db, func(tx *sql.Tx) error {
		// 检查 domain 是否存在
		exists, err := ZoneIsExist(tx, domain)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("domain %s does not exist, please create it first", domain)
		}

		// 查询 records
		rows, err := tx.Query("SELECT id, domain, name, type, value FROM records WHERE domain=?", domain)
		if err != nil {
			return err
		}
		defer func(rows *sql.Rows) {
			_ = rows.Close()
		}(rows)

		// 迭代查询结果
		for rows.Next() {
			var r types.Record
			if err := rows.Scan(&r.Id, &r.Domain, &r.Name, &r.Type, &r.Value); err != nil {
				return err
			}
			records = append(records, r)
		}

		// 检查迭代过程中的错误
		if err := rows.Err(); err != nil {
			return err
		}

		// 如果未找到记录，返回错误
		if len(records) == 0 {
			return fmt.Errorf("no records found for domain %s", domain)
		}

		return nil
	})

	// 返回事务中的错误或结果
	if err != nil {
		return nil, err
	}

	return records, nil
}

// CreateRecord 创建 record
func CreateRecord(db *sql.DB, record types.CreateRecord) error {
	return withTransaction(db, func(tx *sql.Tx) error {
		exists, err := ZoneIsExist(tx, record.Domain)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("domain %s does not exist, please create it first", record.Domain)
		}

		if !validTypes[record.Type] {
			return fmt.Errorf("invalid record type %s", record.Type)
		}

		stmt, err := tx.Prepare("INSERT INTO records(domain, name, type, value) VALUES (?, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer func(stmt *sql.Stmt) {
			_ = stmt.Close()
		}(stmt)
		_, err = stmt.Exec(record.Domain, record.Name, record.Type, record.Value)
		if err != nil {
			var sqliteErr sqlite3.Error
			ok := errors.As(err, &sqliteErr)
			if ok {
				// 检查唯一约束错误
				if errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
					// 进一步验证错误消息确实指示了唯一约束失败
					if strings.Contains(sqliteErr.Error(), "UNIQUE constraint failed") {
						return fmt.Errorf("record already exists")
					}
				}
			}
			return err
		}

		// 获取当前时间戳
		serial := time.Now().Unix()
		stmt, err = tx.Prepare("UPDATE zones SET soa_serial=? WHERE domain=?")
		if err != nil {
			return err
		}
		defer func(stmt *sql.Stmt) {
			_ = stmt.Close()
		}(stmt)
		_, err = stmt.Exec(serial, record.Domain)
		if err != nil {
			return err
		}

		return nil
	})
}

// UpdateRecord 更新 record
func UpdateRecord(db *sql.DB, record types.Record) error {
	return withTransaction(db, func(tx *sql.Tx) error {
		exists, err := RecordIsExist(tx, record.Id)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("record id %d does not exist", record.Id)
		}

		exists, err = ZoneIsExist(tx, record.Domain)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("domain %s does not exist, please create it first", record.Domain)
		}

		if !validTypes[record.Type] {
			return fmt.Errorf("invalid record type %s", record.Type)
		}

		stmt, err := tx.Prepare("UPDATE records SET domain=?, name=?, type=?, value=? WHERE id=?")
		if err != nil {
			return err
		}
		defer func(stmt *sql.Stmt) {
			_ = stmt.Close()
		}(stmt)
		_, err = stmt.Exec(record.Domain, record.Name, record.Type, record.Value, record.Id)
		if err != nil {
			var sqliteErr sqlite3.Error
			if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return fmt.Errorf("record already exists")
			}
			return err
		}

		// 获取当前时间戳
		serial := time.Now().Unix()
		stmt, err = tx.Prepare("UPDATE zones SET soa_serial=? WHERE domain=?")
		if err != nil {
			return err
		}
		defer func(stmt *sql.Stmt) {
			_ = stmt.Close()
		}(stmt)
		_, err = stmt.Exec(serial, record.Domain)
		if err != nil {
			return err
		}

		return nil
	})
}

// DeleteRecord 删除 record
func DeleteRecord(db *sql.DB, id int) error {
	return withTransaction(db, func(tx *sql.Tx) error {
		exists, err := RecordIsExist(tx, id)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("record id %d does not exist", id)
		}

		record, err := GetRecordById(db, id)
		if err != nil {
			return err
		}

		stmt, err := tx.Prepare("DELETE FROM records WHERE id=?")
		if err != nil {
			return err
		}
		defer func(stmt *sql.Stmt) {
			_ = stmt.Close()
		}(stmt)
		_, err = stmt.Exec(id)
		if err != nil {
			return err
		}

		// 获取当前时间戳
		serial := time.Now().Unix()
		stmt, err = tx.Prepare("UPDATE zones SET soa_serial=? WHERE domain=?")
		if err != nil {
			return err
		}
		defer func(stmt *sql.Stmt) {
			_ = stmt.Close()
		}(stmt)
		_, err = stmt.Exec(serial, record.Domain)
		if err != nil {
			return err
		}

		return nil
	})
}

// GetRecordById 根据 record id获取 record
func GetRecordById(db *sql.DB, id int) (record types.Record, err error) {
	err = withTransaction(db, func(tx *sql.Tx) error {
		exists, err := RecordIsExist(tx, id)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("record id %d does not exist", id)
		}

		query := "SELECT domain, name, type, value FROM records WHERE id = ?"
		err = tx.QueryRow(query, id).Scan(&record.Domain, &record.Name, &record.Type, &record.Value)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("no record found for record id %d", id)
			}
			return err
		}

		return nil
	})

	return record, err
}

// GetConfig 获取config
func GetConfig(db *sql.DB, key string) (types.Config, error) {
	var config types.Config
	query := "SELECT value FROM config WHERE key = ?"
	err := db.QueryRow(query, key).Scan(&config.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.Config{Key: key, Value: "no config found for key"}, fmt.Errorf("no config found for key: %s", key)
		}
		return types.Config{Key: key, Value: "no config found for key"}, err
	}
	return types.Config{Key: key, Value: config.Value}, nil
}

// UpdateConfig 更新config
func UpdateConfig(db *sql.DB, config types.Config) error {
	query := "UPDATE config SET value = ? WHERE key = ?"
	result, err := db.Exec(query, config.Value, config.Key)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no config found for key: %s", config.Key)
	}
	return nil
}
