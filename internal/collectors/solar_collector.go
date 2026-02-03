package collectors

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"vpp-go/internal/models"

	"database/sql"
)

// SolarCollector 太陽能數據收集器
type SolarCollector struct {
	DB       *sql.DB
	Model    *models.SolarDataModel
	APIURL   string
	SiteID   string
	Username string
	Password string
}

// NewSolarCollector 創建太陽能數據收集器
func NewSolarCollector(db *sql.DB, apiURL, siteID, username, password string) *SolarCollector {
	return &SolarCollector{
		DB:       db,
		Model:    models.NewSolarDataModel(db),
		APIURL:   apiURL,
		SiteID:   siteID,
		Username: username,
		Password: password,
	}
}

// YihongAPIResponse 義鴻API響應結構
type YihongAPIResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

// FetchData 從義鴻API獲取數據
func (c *SolarCollector) FetchData() (*models.SolarData, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	// 構建請求
	req, err := http.NewRequest("GET", c.APIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("創建請求失敗: %w", err)
	}

	// 添加認證信息（如果需要）
	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	// 發送請求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("請求失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回錯誤狀態碼: %d", resp.StatusCode)
	}

	// 讀取響應
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("讀取響應失敗: %w", err)
	}

	// 解析JSON
	var apiResp YihongAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析JSON失敗: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API返回失敗狀態")
	}

	// 轉換為SolarData結構
	data := c.parseData(apiResp.Data)
	return data, nil
}

// parseData 解析API數據
func (c *SolarCollector) parseData(data map[string]interface{}) *models.SolarData {
	solarData := &models.SolarData{
		SiteID:   c.SiteID,
		DateTime: time.Now(),
	}

	// 解析各個欄位（根據實際API響應格式調整）
	if val, ok := data["daily_generation"].(float64); ok {
		solarData.DailyGeneration = val
	}
	if val, ok := data["solar_radiation"].(float64); ok {
		solarData.SolarRadiation = val
	}
	if val, ok := data["ac_avg_voltage"].(float64); ok {
		solarData.ACAverageVoltage = val
	}
	if val, ok := data["ac_total_power"].(float64); ok {
		solarData.ACTotalPower = val
	}
	if val, ok := data["ac_total_current"].(float64); ok {
		solarData.ACTotalCurrent = val
	}
	if val, ok := data["dc_avg_voltage"].(float64); ok {
		solarData.DCAverageVoltage = val
	}
	if val, ok := data["dc_total_power"].(float64); ok {
		solarData.DCTotalPower = val
	}
	if val, ok := data["dc_total_current"].(float64); ok {
		solarData.DCTotalCurrent = val
	}
	if val, ok := data["module_temperature"].(float64); ok {
		solarData.ModuleTemperature = val
	}
	if val, ok := data["total_accumulated_generation"].(float64); ok {
		solarData.TotalAccumulatedGeneration = val
	}
	if val, ok := data["co2_reduction"].(float64); ok {
		solarData.CO2Reduction = val
	}

	return solarData
}

// SaveToDatabase 保存數據到資料庫
func (c *SolarCollector) SaveToDatabase(data *models.SolarData) error {
	return c.Model.Insert(data)
}

// CollectAndSave 收集並保存數據
func (c *SolarCollector) CollectAndSave() error {
	log.Printf("開始收集太陽能數據 - 場站: %s\n", c.SiteID)

	data, err := c.FetchData()
	if err != nil {
		return fmt.Errorf("獲取數據失敗: %w", err)
	}

	if err := c.SaveToDatabase(data); err != nil {
		return fmt.Errorf("保存數據失敗: %w", err)
	}

	log.Printf("太陽能數據收集成功 - 場站: %s, 時間: %s\n", c.SiteID, data.DateTime.Format("2006-01-02 15:04:05"))
	return nil
}

// StartSchedule 啟動定時收集（每15分鐘）
func (c *SolarCollector) StartSchedule() {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	// 立即執行一次
	if err := c.CollectAndSave(); err != nil {
		log.Printf("太陽能數據收集錯誤: %v\n", err)
	}

	// 定時執行
	for range ticker.C {
		if err := c.CollectAndSave(); err != nil {
			log.Printf("太陽能數據收集錯誤: %v\n", err)
		}
	}
}
