package database

import (
	"database/sql"
	"fmt"
	"log"
	"vpp-go/internal/config"

	_ "github.com/lib/pq"
)

// InitDB 初始化資料庫連接
func InitDB(cfg *config.Config) (*sql.DB, error) {
	dsn := cfg.GetDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("無法打開資料庫連接: %w", err)
	}

	// 測試連接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("無法連接到資料庫: %w", err)
	}

	// 設置連接池
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("資料庫連接成功")
	return db, nil
}

// ExecuteQuery 執行查詢並返回結果
func ExecuteQuery(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("查詢執行失敗: %w", err)
	}
	return rows, nil
}

// ExecuteQueryRow 執行查詢並返回單行結果
func ExecuteQueryRow(db *sql.DB, query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}

// ExecuteInsert 執行插入操作
func ExecuteInsert(db *sql.DB, query string, args ...interface{}) (int64, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("插入操作失敗: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		// PostgreSQL 不支持 LastInsertId，返回受影響的行數
		rowsAffected, _ := result.RowsAffected()
		return rowsAffected, nil
	}

	return id, nil
}

// ExecuteBatch 批次執行插入操作
func ExecuteBatch(db *sql.DB, query string, dataList [][]interface{}) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("無法開始事務: %w", err)
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("無法準備語句: %w", err)
	}
	defer stmt.Close()

	for _, data := range dataList {
		if _, err := stmt.Exec(data...); err != nil {
			tx.Rollback()
			return fmt.Errorf("批次執行失敗: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("事務提交失敗: %w", err)
	}

	return nil
}
