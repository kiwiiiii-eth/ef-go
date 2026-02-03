package models

import (
	"database/sql"
	"time"
)

// LoadData 負載數據模型
type LoadData struct {
	ID        int       `json:"id"`
	SiteID    string    `json:"site_id"`
	DateTime  time.Time `json:"datetime"`
	LoadValue float64   `json:"load_value"`
}

// LoadDataModel 負載數據模型操作
type LoadDataModel struct {
	DB *sql.DB
}

// NewLoadDataModel 創建負載數據模型
func NewLoadDataModel(db *sql.DB) *LoadDataModel {
	return &LoadDataModel{DB: db}
}

// GetLatest 獲取最新的負載數據
func (m *LoadDataModel) GetLatest(siteID string) (*LoadData, error) {
	query := `
		SELECT id, site_id, datetime, load_value
		FROM load_data
		WHERE site_id = $1
		ORDER BY datetime DESC
		LIMIT 1
	`

	data := &LoadData{}
	err := m.DB.QueryRow(query, siteID).Scan(
		&data.ID, &data.SiteID, &data.DateTime, &data.LoadValue,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetAllLatest 獲取所有場站最新的負載數據
func (m *LoadDataModel) GetAllLatest() ([]LoadData, error) {
	query := `
		SELECT DISTINCT ON (site_id)
		       id, site_id, datetime, load_value
		FROM load_data
		ORDER BY site_id, datetime DESC
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []LoadData
	for rows.Next() {
		var data LoadData
		err := rows.Scan(&data.ID, &data.SiteID, &data.DateTime, &data.LoadValue)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, data)
	}

	return dataList, nil
}

// GetHistory 獲取歷史數據
func (m *LoadDataModel) GetHistory(siteID string, startDate, endDate time.Time, limit int) ([]LoadData, error) {
	query := `
		SELECT id, site_id, datetime, load_value
		FROM load_data
		WHERE site_id = $1 AND datetime BETWEEN $2 AND $3
		ORDER BY datetime DESC
		LIMIT $4
	`

	rows, err := m.DB.Query(query, siteID, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []LoadData
	for rows.Next() {
		var data LoadData
		err := rows.Scan(&data.ID, &data.SiteID, &data.DateTime, &data.LoadValue)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, data)
	}

	return dataList, nil
}

// Insert 插入負載數據
func (m *LoadDataModel) Insert(data *LoadData) error {
	query := `
		INSERT INTO load_data (site_id, datetime, load_value)
		VALUES ($1, $2, $3)
		ON CONFLICT (site_id, datetime) DO UPDATE SET
			load_value = EXCLUDED.load_value
	`

	_, err := m.DB.Exec(query, data.SiteID, data.DateTime, data.LoadValue)
	return err
}
