package models

import (
	"database/sql"
	"time"
)

// TaipowerReserveData 台電備轉資料模型
type TaipowerReserveData struct {
	ID              int       `json:"id"`
	TranDate        time.Time `json:"tran_date"`
	TranHour        int       `json:"tran_hour"`
	SRBid           float64   `json:"sr_bid"`
	SRBidQSE        float64   `json:"sr_bid_qse"`
	SRBidNonTrade   float64   `json:"sr_bid_nontrade"`
	SRPrice         float64   `json:"sr_price"`
	SRPerfPrice1    float64   `json:"sr_perf_price_1"`
	SRPerfPrice2    float64   `json:"sr_perf_price_2"`
	SRPerfPrice3    float64   `json:"sr_perf_price_3"`
	SUPBid          float64   `json:"sup_bid"`
	SUPBidQSE       float64   `json:"sup_bid_qse"`
	SUPBidNonTrade  float64   `json:"sup_bid_nontrade"`
	SUPPrice        float64   `json:"sup_price"`
}

// TaipowerReserveModel 台電備轉資料模型操作
type TaipowerReserveModel struct {
	DB *sql.DB
}

// NewTaipowerReserveModel 創建台電備轉資料模型
func NewTaipowerReserveModel(db *sql.DB) *TaipowerReserveModel {
	return &TaipowerReserveModel{DB: db}
}

// GetLatest 獲取最新一天的備轉資料
func (m *TaipowerReserveModel) GetLatest() ([]TaipowerReserveData, error) {
	query := `
		SELECT id, tran_date, tran_hour, sr_bid, sr_bid_qse, sr_bid_nontrade,
		       sr_price, sr_perf_price_1, sr_perf_price_2, sr_perf_price_3,
		       sup_bid, sup_bid_qse, sup_bid_nontrade, sup_price
		FROM taipower_reserve_data
		WHERE tran_date = (SELECT MAX(tran_date) FROM taipower_reserve_data)
		ORDER BY tran_hour
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []TaipowerReserveData
	for rows.Next() {
		var data TaipowerReserveData
		err := rows.Scan(
			&data.ID, &data.TranDate, &data.TranHour,
			&data.SRBid, &data.SRBidQSE, &data.SRBidNonTrade,
			&data.SRPrice, &data.SRPerfPrice1, &data.SRPerfPrice2, &data.SRPerfPrice3,
			&data.SUPBid, &data.SUPBidQSE, &data.SUPBidNonTrade, &data.SUPPrice,
		)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, data)
	}

	return dataList, nil
}

// GetByDate 獲取特定日期的備轉資料
func (m *TaipowerReserveModel) GetByDate(date time.Time) ([]TaipowerReserveData, error) {
	query := `
		SELECT id, tran_date, tran_hour, sr_bid, sr_bid_qse, sr_bid_nontrade,
		       sr_price, sr_perf_price_1, sr_perf_price_2, sr_perf_price_3,
		       sup_bid, sup_bid_qse, sup_bid_nontrade, sup_price
		FROM taipower_reserve_data
		WHERE tran_date = $1
		ORDER BY tran_hour
	`

	rows, err := m.DB.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []TaipowerReserveData
	for rows.Next() {
		var data TaipowerReserveData
		err := rows.Scan(
			&data.ID, &data.TranDate, &data.TranHour,
			&data.SRBid, &data.SRBidQSE, &data.SRBidNonTrade,
			&data.SRPrice, &data.SRPerfPrice1, &data.SRPerfPrice2, &data.SRPerfPrice3,
			&data.SUPBid, &data.SUPBidQSE, &data.SUPBidNonTrade, &data.SUPPrice,
		)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, data)
	}

	return dataList, nil
}

// GetHistory 獲取歷史備轉資料
func (m *TaipowerReserveModel) GetHistory(startDate, endDate time.Time, limit int) ([]TaipowerReserveData, error) {
	query := `
		SELECT id, tran_date, tran_hour, sr_bid, sr_bid_qse, sr_bid_nontrade,
		       sr_price, sr_perf_price_1, sr_perf_price_2, sr_perf_price_3,
		       sup_bid, sup_bid_qse, sup_bid_nontrade, sup_price
		FROM taipower_reserve_data
		WHERE tran_date BETWEEN $1 AND $2
		ORDER BY tran_date DESC, tran_hour DESC
		LIMIT $3
	`

	rows, err := m.DB.Query(query, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []TaipowerReserveData
	for rows.Next() {
		var data TaipowerReserveData
		err := rows.Scan(
			&data.ID, &data.TranDate, &data.TranHour,
			&data.SRBid, &data.SRBidQSE, &data.SRBidNonTrade,
			&data.SRPrice, &data.SRPerfPrice1, &data.SRPerfPrice2, &data.SRPerfPrice3,
			&data.SUPBid, &data.SUPBidQSE, &data.SUPBidNonTrade, &data.SUPPrice,
		)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, data)
	}

	return dataList, nil
}

// GetByHour 獲取特定時段的備轉資料
func (m *TaipowerReserveModel) GetByHour(date time.Time, hour int) (*TaipowerReserveData, error) {
	query := `
		SELECT id, tran_date, tran_hour, sr_bid, sr_bid_qse, sr_bid_nontrade,
		       sr_price, sr_perf_price_1, sr_perf_price_2, sr_perf_price_3,
		       sup_bid, sup_bid_qse, sup_bid_nontrade, sup_price
		FROM taipower_reserve_data
		WHERE tran_date = $1 AND tran_hour = $2
	`

	data := &TaipowerReserveData{}
	err := m.DB.QueryRow(query, date, hour).Scan(
		&data.ID, &data.TranDate, &data.TranHour,
		&data.SRBid, &data.SRBidQSE, &data.SRBidNonTrade,
		&data.SRPrice, &data.SRPerfPrice1, &data.SRPerfPrice2, &data.SRPerfPrice3,
		&data.SUPBid, &data.SUPBidQSE, &data.SUPBidNonTrade, &data.SUPPrice,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetStatistics 獲取統計資訊
func (m *TaipowerReserveModel) GetStatistics(date time.Time) (map[string]interface{}, error) {
	query := `
		SELECT
			AVG(sr_bid) as sr_bid_avg,
			MAX(sr_bid) as sr_bid_max,
			MIN(sr_bid) as sr_bid_min,
			AVG(sr_price) as sr_price_avg,
			MAX(sr_price) as sr_price_max,
			MIN(sr_price) as sr_price_min,
			AVG(sup_bid) as sup_bid_avg,
			MAX(sup_bid) as sup_bid_max,
			MIN(sup_bid) as sup_bid_min,
			AVG(sup_price) as sup_price_avg,
			MAX(sup_price) as sup_price_max,
			MIN(sup_price) as sup_price_min
		FROM taipower_reserve_data
		WHERE tran_date = $1
	`

	var stats map[string]interface{} = make(map[string]interface{})
	var srBidAvg, srBidMax, srBidMin, srPriceAvg, srPriceMax, srPriceMin sql.NullFloat64
	var supBidAvg, supBidMax, supBidMin, supPriceAvg, supPriceMax, supPriceMin sql.NullFloat64

	err := m.DB.QueryRow(query, date).Scan(
		&srBidAvg, &srBidMax, &srBidMin,
		&srPriceAvg, &srPriceMax, &srPriceMin,
		&supBidAvg, &supBidMax, &supBidMin,
		&supPriceAvg, &supPriceMax, &supPriceMin,
	)

	if err != nil {
		return nil, err
	}

	stats["sr_bid_avg"] = getFloatValue(srBidAvg)
	stats["sr_bid_max"] = getFloatValue(srBidMax)
	stats["sr_bid_min"] = getFloatValue(srBidMin)
	stats["sr_price_avg"] = getFloatValue(srPriceAvg)
	stats["sr_price_max"] = getFloatValue(srPriceMax)
	stats["sr_price_min"] = getFloatValue(srPriceMin)
	stats["sup_bid_avg"] = getFloatValue(supBidAvg)
	stats["sup_bid_max"] = getFloatValue(supBidMax)
	stats["sup_bid_min"] = getFloatValue(supBidMin)
	stats["sup_price_avg"] = getFloatValue(supPriceAvg)
	stats["sup_price_max"] = getFloatValue(supPriceMax)
	stats["sup_price_min"] = getFloatValue(supPriceMin)

	return stats, nil
}

// Insert 插入備轉資料
func (m *TaipowerReserveModel) Insert(data *TaipowerReserveData) error {
	query := `
		INSERT INTO taipower_reserve_data (
			tran_date, tran_hour, sr_bid, sr_bid_qse, sr_bid_nontrade,
			sr_price, sr_perf_price_1, sr_perf_price_2, sr_perf_price_3,
			sup_bid, sup_bid_qse, sup_bid_nontrade, sup_price
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (tran_date, tran_hour) DO UPDATE SET
			sr_bid = EXCLUDED.sr_bid,
			sr_bid_qse = EXCLUDED.sr_bid_qse,
			sr_bid_nontrade = EXCLUDED.sr_bid_nontrade,
			sr_price = EXCLUDED.sr_price,
			sr_perf_price_1 = EXCLUDED.sr_perf_price_1,
			sr_perf_price_2 = EXCLUDED.sr_perf_price_2,
			sr_perf_price_3 = EXCLUDED.sr_perf_price_3,
			sup_bid = EXCLUDED.sup_bid,
			sup_bid_qse = EXCLUDED.sup_bid_qse,
			sup_bid_nontrade = EXCLUDED.sup_bid_nontrade,
			sup_price = EXCLUDED.sup_price
	`

	_, err := m.DB.Exec(query,
		data.TranDate, data.TranHour,
		data.SRBid, data.SRBidQSE, data.SRBidNonTrade,
		data.SRPrice, data.SRPerfPrice1, data.SRPerfPrice2, data.SRPerfPrice3,
		data.SUPBid, data.SUPBidQSE, data.SUPBidNonTrade, data.SUPPrice,
	)

	return err
}

// getFloatValue 輔助函數：從 sql.NullFloat64 獲取值
func getFloatValue(nf sql.NullFloat64) float64 {
	if nf.Valid {
		return nf.Float64
	}
	return 0
}
