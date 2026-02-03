package collectors

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"vpp-go/internal/models"

	"database/sql"
)

// TaipowerCollector 台電備轉資料收集器
type TaipowerCollector struct {
	DB     *sql.DB
	Model  *models.TaipowerReserveModel
	BaseURL string
}

// NewTaipowerCollector 創建台電備轉資料收集器
func NewTaipowerCollector(db *sql.DB, baseURL string) *TaipowerCollector {
	return &TaipowerCollector{
		DB:      db,
		Model:   models.NewTaipowerReserveModel(db),
		BaseURL: baseURL,
	}
}

// FetchData 爬取台電網站數據
func (c *TaipowerCollector) FetchData(date time.Time) ([]models.TaipowerReserveData, error) {
	dateStr := date.Format("20060102") // YYYYMMDD格式

	// 構建URL（根據實際台電網站調整）
	url := fmt.Sprintf("%s/reserve_data?date=%s", c.BaseURL, dateStr)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("請求失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("網站返回錯誤狀態碼: %d", resp.StatusCode)
	}

	// 讀取HTML內容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("讀取響應失敗: %w", err)
	}

	// 解析HTML並提取數據
	dataList := c.parseHTML(string(body), date)
	return dataList, nil
}

// parseHTML 解析HTML並提取數據
func (c *TaipowerCollector) parseHTML(html string, date time.Time) []models.TaipowerReserveData {
	var dataList []models.TaipowerReserveData

	// 這裡需要根據實際的台電網站HTML結構來解析
	// 以下是一個簡化的示例，實際使用時需要根據網站結構調整

	// 使用正則表達式提取表格數據
	tableRegex := regexp.MustCompile(`<tr[^>]*>(.*?)</tr>`)
	rowRegex := regexp.MustCompile(`<td[^>]*>(.*?)</td>`)

	tables := tableRegex.FindAllStringSubmatch(html, -1)
	for _, table := range tables {
		cells := rowRegex.FindAllStringSubmatch(table[1], -1)
		if len(cells) >= 14 { // 確保有足夠的欄位
			data := models.TaipowerReserveData{
				TranDate: date,
			}

			// 解析每個欄位（根據實際網站結構調整索引）
			data.TranHour = c.parseInt(c.cleanHTML(cells[0][1]))
			data.SRBid = c.parseFloat(c.cleanHTML(cells[1][1]))
			data.SRBidQSE = c.parseFloat(c.cleanHTML(cells[2][1]))
			data.SRBidNonTrade = c.parseFloat(c.cleanHTML(cells[3][1]))
			data.SRPrice = c.parseFloat(c.cleanHTML(cells[4][1]))
			data.SRPerfPrice1 = c.parseFloat(c.cleanHTML(cells[5][1]))
			data.SRPerfPrice2 = c.parseFloat(c.cleanHTML(cells[6][1]))
			data.SRPerfPrice3 = c.parseFloat(c.cleanHTML(cells[7][1]))
			data.SUPBid = c.parseFloat(c.cleanHTML(cells[8][1]))
			data.SUPBidQSE = c.parseFloat(c.cleanHTML(cells[9][1]))
			data.SUPBidNonTrade = c.parseFloat(c.cleanHTML(cells[10][1]))
			data.SUPPrice = c.parseFloat(c.cleanHTML(cells[11][1]))

			dataList = append(dataList, data)
		}
	}

	return dataList
}

// cleanHTML 清理HTML標籤
func (c *TaipowerCollector) cleanHTML(s string) string {
	// 移除HTML標籤
	re := regexp.MustCompile(`<[^>]*>`)
	s = re.ReplaceAllString(s, "")
	// 移除空白字符
	s = strings.TrimSpace(s)
	return s
}

// parseInt 解析整數
func (c *TaipowerCollector) parseInt(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

// parseFloat 解析浮點數
func (c *TaipowerCollector) parseFloat(s string) float64 {
	// 移除逗號
	s = strings.ReplaceAll(s, ",", "")
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

// SaveToDatabase 批次保存數據到資料庫
func (c *TaipowerCollector) SaveToDatabase(dataList []models.TaipowerReserveData) error {
	for _, data := range dataList {
		if err := c.Model.Insert(&data); err != nil {
			return fmt.Errorf("保存數據失敗: %w", err)
		}
	}
	return nil
}

// CollectAndSave 收集並保存數據
func (c *TaipowerCollector) CollectAndSave(date time.Time) error {
	log.Printf("開始收集台電備轉資料 - 日期: %s\n", date.Format("2006-01-02"))

	dataList, err := c.FetchData(date)
	if err != nil {
		return fmt.Errorf("獲取數據失敗: %w", err)
	}

	if len(dataList) == 0 {
		return fmt.Errorf("沒有找到數據")
	}

	if err := c.SaveToDatabase(dataList); err != nil {
		return fmt.Errorf("保存數據失敗: %w", err)
	}

	log.Printf("台電備轉資料收集成功 - 日期: %s, 筆數: %d\n", date.Format("2006-01-02"), len(dataList))
	return nil
}

// StartSchedule 啟動定時收集（每天凌晨2點）
func (c *TaipowerCollector) StartSchedule() {
	// 計算下次執行時間（凌晨2點）
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
	if now.After(next) {
		next = next.Add(24 * time.Hour)
	}

	// 等待到下次執行時間
	time.Sleep(time.Until(next))

	// 執行任務
	yesterday := time.Now().AddDate(0, 0, -1)
	if err := c.CollectAndSave(yesterday); err != nil {
		log.Printf("台電備轉資料收集錯誤: %v\n", err)
	}

	// 每24小時執行一次
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		yesterday := time.Now().AddDate(0, 0, -1)
		if err := c.CollectAndSave(yesterday); err != nil {
			log.Printf("台電備轉資料收集錯誤: %v\n", err)
		}
	}
}
