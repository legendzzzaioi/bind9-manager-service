package model

import (
	"bind9-manager-service/internal/types"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 生成 bind9 所有配置
func GenerateAllZoneFiles(db *sql.DB, bindpath string) error {
	// 确保 bindpath 目录存在
	if err := os.MkdirAll(bindpath, 0755); err != nil {
		return fmt.Errorf("failed to create bindpath directory: %w", err)
	}

	// 获取所有 zone
	zones, err := GetZones(db)
	if err != nil {
		if strings.Contains(err.Error(), "no zones found") {
			// 如果没有 zone，则清空 named.conf.local 文件
			namedLocalConfFile, err := os.OpenFile(filepath.Join(bindpath, "named.conf.local"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return fmt.Errorf("failed to open named.conf.local: %w", err)
			}
			defer namedLocalConfFile.Close()
			return nil
		}
		return fmt.Errorf("failed to get zones: %w", err)
	}

	// 打开 named.conf.local 文件进行重写操作
	namedLocalConfFile, err := os.OpenFile(filepath.Join(bindpath, "named.conf.local"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open named.conf.local: %w", err)
	}
	defer namedLocalConfFile.Close()

	// 迭代所有 zones
	var fileContent string
	for _, zone := range zones {
		// 生成并写入 zone 文件
		if err := GenerateZoneFile(db, bindpath, zone); err != nil {
			return fmt.Errorf("failed to generate zone file for domain %s: %w", zone.Domain, err)
		}

		// 构建named.conf.local文件内容
		fileName := filepath.Join(bindpath, fmt.Sprintf("db-%s", zone.Domain))
		zoneConfig := fmt.Sprintf(`zone "%s" {type primary; file "%s";};`+"\n", zone.Domain, fileName)
		fileContent += zoneConfig
	}

	// 写入 named.conf.local 文件
	if _, err := namedLocalConfFile.WriteString(fileContent); err != nil {
		return fmt.Errorf("failed to write zone configurations to named.conf.local: %w", err)
	}

	return nil
}

// 生成 named.conf.local 文件
func GenerateNamedLocalConf(db *sql.DB, bindpath string) error {
	// 确保 bindpath 目录存在
	if err := os.MkdirAll(bindpath, 0755); err != nil {
		return fmt.Errorf("failed to create bindpath directory: %w", err)
	}

	// 获取所有 zone
	zones, err := GetZones(db)
	if err != nil {
		if strings.Contains(err.Error(), "no zones found") {
			// 如果没有 zone，则清空 named.conf.local 文件
			namedLocalConfFile, err := os.OpenFile(filepath.Join(bindpath, "named.conf.local"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return fmt.Errorf("failed to open named.conf.local: %w", err)
			}
			defer namedLocalConfFile.Close()
			return nil
		}
		return fmt.Errorf("failed to get zones: %w", err)
	}

	// 打开 named.conf.local 文件进行重写操作
	namedLocalConfFile, err := os.OpenFile(filepath.Join(bindpath, "named.conf.local"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open named.conf.local: %w", err)
	}
	defer namedLocalConfFile.Close()

	// 迭代所有 zones，构建named.conf.local文件内容
	var fileContent string
	for _, zone := range zones {
		fileName := filepath.Join(bindpath, fmt.Sprintf("db-%s", zone.Domain))
		zoneConfig := fmt.Sprintf(`zone "%s" {type primary; file "%s";};`+"\n", zone.Domain, fileName)
		fileContent += zoneConfig
	}

	// 写入文件
	if _, err := namedLocalConfFile.WriteString(fileContent); err != nil {
		return fmt.Errorf("failed to write zone configurations to named.conf.local: %w", err)
	}

	return nil
}

// 生成并写入 zone 文件
func GenerateZoneFile(db *sql.DB, bindpath string, zone types.Zone) error {
	fileName := filepath.Join(bindpath, fmt.Sprintf("db-%s", zone.Domain))
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create zone file for domain %s: %w", zone.Domain, err)
	}
	defer file.Close()

	// 写入 zone 文件头部信息
	if err := WriteZoneHeader(file, zone); err != nil {
		return fmt.Errorf("failed to write SOA record for domain %s: %w", zone.Domain, err)
	}

	// 查询和写入与当前 zone 相关的 records
	if err := WriteZoneRecords(db, file, zone.Domain); err != nil {
		return fmt.Errorf("failed to write records to zone file for domain %s: %w", zone.Domain, err)
	}

	return nil
}

// 根据 domain 生成并写入 zone 文件
func GenerateZoneFileByDomain(db *sql.DB, bindpath string, domain string) error {
	fileName := filepath.Join(bindpath, fmt.Sprintf("db-%s", domain))
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create zone file for domain %s: %w", domain, err)
	}
	defer file.Close()

	// 获取 zone
	zone, err := GetZoneByDomain(db, domain)
	if err != nil {
		return fmt.Errorf("failed to get zone for domain %s: %w", domain, err)
	}

	// 写入 zone 文件头部信息
	if err := WriteZoneHeader(file, zone); err != nil {
		return fmt.Errorf("failed to write SOA record for domain %s: %w", domain, err)
	}

	// 查询和写入与当前 zone 相关的 records
	if err := WriteZoneRecords(db, file, domain); err != nil {
		return fmt.Errorf("failed to write records to zone file for domain %s: %w", domain, err)
	}

	return nil
}

// 写入 zone 文件头部信息
func WriteZoneHeader(file *os.File, zone types.Zone) error {
	// SOA 记录
	soa := fmt.Sprintf("$TTL %d\n@ IN SOA %s. %s. (\n\t%d ; Serial\n\t%d ; Refresh\n\t%d ; Retry\n\t%d ; Expire\n\t%d ) ; Negative Cache TTL\n",
		zone.Ttl, zone.PrimaryNameServer, strings.ReplaceAll(zone.MailAddress, "@", "."), zone.Serial, zone.Refresh, zone.Retry, zone.Expire, zone.CacheTtl)

	// NS 记录
	ns := fmt.Sprintf("@ IN NS %s.\n", zone.PrimaryNameServer)

	// 写入文件
	_, err := file.WriteString(soa)
	if err != nil {
		return err
	}
	_, err = file.WriteString(ns)
	return err
}

// 写入 zone 文件中的记录信息
func WriteZoneRecords(db *sql.DB, file *os.File, domain string) error {
	query := "SELECT domain, name, type, value FROM records WHERE domain = ?"
	rows, err := db.Query(query, domain)
	if err != nil {
		return fmt.Errorf("failed to query records for domain %s: %w", domain, err)
	}
	defer rows.Close()

	for rows.Next() {
		var record types.CreateRecord
		if err := rows.Scan(&record.Domain, &record.Name, &record.Type, &record.Value); err != nil {
			return fmt.Errorf("failed to scan record: %w", err)
		}
		recordLine := fmt.Sprintf("%s IN %s %s\n", record.Name, record.Type, record.Value)
		if _, err := file.WriteString(recordLine); err != nil {
			return fmt.Errorf("failed to write record to zone file for domain %s: %w", record.Domain, err)
		}
	}

	return rows.Err()
}

// 根据 domain 删除 zone 文件
func DeleteZoneFileByDomain(bindpath string, domain string) error {
	fileName := filepath.Join(bindpath, fmt.Sprintf("db-%s", domain))

	// 检查文件是否存在
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", fileName)
	}

	// 尝试删除文件
	if err := os.Remove(fileName); err != nil {
		return fmt.Errorf("failed to delete file %s: %w", fileName, err)
	}

	return nil
}
