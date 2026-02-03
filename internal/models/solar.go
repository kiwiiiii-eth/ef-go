package models

import (
	"database/sql"
	"time"
)

// SolarData 太陽能數據模型
type SolarData struct {
	ID                         int       `json:"id"`
	SiteID                     string    `json:"site_id"`
	DateTime                   time.Time `json:"datetime"`
	DailyGeneration            float64   `json:"daily_generation"`
	SolarRadiation             float64   `json:"solar_radiation"`
	ACAverageVoltage           float64   `json:"ac_avg_voltage"`
	ACTotalPower               float64   `json:"ac_total_power"`
	ACTotalCurrent             float64   `json:"ac_total_current"`
	DCAverageVoltage           float64   `json:"dc_avg_voltage"`
	DCTotalPower               float64   `json:"dc_total_power"`
	DCTotalCurrent             float64   `json:"dc_total_current"`
	ModuleTemperature          float64   `json:"module_temperature"`
	TotalAccumulatedGeneration float64   `json:"total_accumulated_generation"`
	CO2Reduction               float64   `json:"co2_reduction"`
}

// SolarDataModel 太陽能數據模型操作
type SolarDataModel struct {
	DB *sql.DB
}

// NewSolarDataModel 創建太陽能數據模型
func NewSolarDataModel(db *sql.DB) *SolarDataModel {
	return &SolarDataModel{DB: db}
}

// GetLatest 獲取最新的太陽能數據
func (m *SolarDataModel) GetLatest(siteID string) (*SolarData, error) {
	query := `
		SELECT id, site_id, datetime, daily_generation, solar_radiation,
		       ac_avg_voltage, ac_total_power, ac_total_current,
		       dc_avg_voltage, dc_total_power, dc_total_current,
		       module_temperature, total_accumulated_generation, co2_reduction
		FROM solar_data
		WHERE site_id = $1
		ORDER BY datetime DESC
		LIMIT 1
	`

	data := &SolarData{}
	err := m.DB.QueryRow(query, siteID).Scan(
		&data.ID, &data.SiteID, &data.DateTime, &data.DailyGeneration,
		&data.SolarRadiation, &data.ACAverageVoltage, &data.ACTotalPower,
		&data.ACTotalCurrent, &data.DCAverageVoltage, &data.DCTotalPower,
		&data.DCTotalCurrent, &data.ModuleTemperature,
		&data.TotalAccumulatedGeneration, &data.CO2Reduction,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetAllLatest 獲取所有場站最新的太陽能數據
func (m *SolarDataModel) GetAllLatest() ([]SolarData, error) {
	query := `
		SELECT DISTINCT ON (site_id)
		       id, site_id, datetime, daily_generation, solar_radiation,
		       ac_avg_voltage, ac_total_power, ac_total_current,
		       dc_avg_voltage, dc_total_power, dc_total_current,
		       module_temperature, total_accumulated_generation, co2_reduction
		FROM solar_data
		ORDER BY site_id, datetime DESC
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []SolarData
	for rows.Next() {
		var data SolarData
		err := rows.Scan(
			&data.ID, &data.SiteID, &data.DateTime, &data.DailyGeneration,
			&data.SolarRadiation, &data.ACAverageVoltage, &data.ACTotalPower,
			&data.ACTotalCurrent, &data.DCAverageVoltage, &data.DCTotalPower,
			&data.DCTotalCurrent, &data.ModuleTemperature,
			&data.TotalAccumulatedGeneration, &data.CO2Reduction,
		)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, data)
	}

	return dataList, nil
}

// GetHistory 獲取歷史數據
func (m *SolarDataModel) GetHistory(siteID string, startDate, endDate time.Time, limit int) ([]SolarData, error) {
	query := `
		SELECT id, site_id, datetime, daily_generation, solar_radiation,
		       ac_avg_voltage, ac_total_power, ac_total_current,
		       dc_avg_voltage, dc_total_power, dc_total_current,
		       module_temperature, total_accumulated_generation, co2_reduction
		FROM solar_data
		WHERE site_id = $1 AND datetime BETWEEN $2 AND $3
		ORDER BY datetime DESC
		LIMIT $4
	`

	rows, err := m.DB.Query(query, siteID, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []SolarData
	for rows.Next() {
		var data SolarData
		err := rows.Scan(
			&data.ID, &data.SiteID, &data.DateTime, &data.DailyGeneration,
			&data.SolarRadiation, &data.ACAverageVoltage, &data.ACTotalPower,
			&data.ACTotalCurrent, &data.DCAverageVoltage, &data.DCTotalPower,
			&data.DCTotalCurrent, &data.ModuleTemperature,
			&data.TotalAccumulatedGeneration, &data.CO2Reduction,
		)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, data)
	}

	return dataList, nil
}

// Insert 插入太陽能數據
func (m *SolarDataModel) Insert(data *SolarData) error {
	query := `
		INSERT INTO solar_data (
			site_id, datetime, daily_generation, solar_radiation,
			ac_avg_voltage, ac_total_power, ac_total_current,
			dc_avg_voltage, dc_total_power, dc_total_current,
			module_temperature, total_accumulated_generation, co2_reduction
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (site_id, datetime) DO UPDATE SET
			daily_generation = EXCLUDED.daily_generation,
			solar_radiation = EXCLUDED.solar_radiation,
			ac_avg_voltage = EXCLUDED.ac_avg_voltage,
			ac_total_power = EXCLUDED.ac_total_power,
			ac_total_current = EXCLUDED.ac_total_current,
			dc_avg_voltage = EXCLUDED.dc_avg_voltage,
			dc_total_power = EXCLUDED.dc_total_power,
			dc_total_current = EXCLUDED.dc_total_current,
			module_temperature = EXCLUDED.module_temperature,
			total_accumulated_generation = EXCLUDED.total_accumulated_generation,
			co2_reduction = EXCLUDED.co2_reduction
	`

	_, err := m.DB.Exec(query,
		data.SiteID, data.DateTime, data.DailyGeneration, data.SolarRadiation,
		data.ACAverageVoltage, data.ACTotalPower, data.ACTotalCurrent,
		data.DCAverageVoltage, data.DCTotalPower, data.DCTotalCurrent,
		data.ModuleTemperature, data.TotalAccumulatedGeneration, data.CO2Reduction,
	)

	return err
}
