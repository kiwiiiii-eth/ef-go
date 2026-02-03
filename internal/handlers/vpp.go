package handlers

import (
	"net/http"
	"strconv"
	"time"
	"vpp-go/internal/config"

	"github.com/gin-gonic/gin"
)

// GetAllRealtimeData 獲取所有場站即時數據
func (h *Handler) GetAllRealtimeData(c *gin.Context) {
	solarData, err := h.SolarModel.GetAllLatest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	loadData, err := h.LoadModel.GetAllLatest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"solar": solarData,
		"load":  loadData,
	})
}

// GetSiteRealtimeData 獲取特定場站即時數據
func (h *Handler) GetSiteRealtimeData(c *gin.Context) {
	siteID := c.Param("site_id")

	if !config.IsValidSite(siteID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的場站ID"})
		return
	}

	solarData, err := h.SolarModel.GetLatest(siteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	loadData, err := h.LoadModel.GetLatest(siteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"site_id": siteID,
		"solar":   solarData,
		"load":    loadData,
	})
}

// GetLatestSolarData 獲取最新太陽能數據
func (h *Handler) GetLatestSolarData(c *gin.Context) {
	siteID := c.Query("site_id")

	if siteID != "" {
		// 獲取特定場站
		if !config.IsValidSite(siteID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無效的場站ID"})
			return
		}

		data, err := h.SolarModel.GetLatest(siteID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if data == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "找不到數據"})
			return
		}

		c.JSON(http.StatusOK, data)
	} else {
		// 獲取所有場站
		dataList, err := h.SolarModel.GetAllLatest()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, dataList)
	}
}

// GetSolarHistory 獲取太陽能歷史數據
func (h *Handler) GetSolarHistory(c *gin.Context) {
	siteID := c.Query("site_id")
	if siteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少場站ID參數"})
		return
	}

	if !config.IsValidSite(siteID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的場站ID"})
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	limitStr := c.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的limit參數"})
		return
	}

	// 解析日期
	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無效的開始日期格式"})
			return
		}
	} else {
		startDate = time.Now().AddDate(0, 0, -30) // 預設30天前
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無效的結束日期格式"})
			return
		}
	} else {
		endDate = time.Now()
	}

	dataList, err := h.SolarModel.GetHistory(siteID, startDate, endDate, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"site_id":    siteID,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"count":      len(dataList),
		"data":       dataList,
	})
}

// GetLatestLoadData 獲取最新負載數據
func (h *Handler) GetLatestLoadData(c *gin.Context) {
	siteID := c.Query("site_id")

	if siteID != "" {
		// 獲取特定場站
		if !config.IsValidSite(siteID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無效的場站ID"})
			return
		}

		data, err := h.LoadModel.GetLatest(siteID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if data == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "找不到數據"})
			return
		}

		c.JSON(http.StatusOK, data)
	} else {
		// 獲取所有場站
		dataList, err := h.LoadModel.GetAllLatest()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, dataList)
	}
}

// GetLoadHistory 獲取負載歷史數據
func (h *Handler) GetLoadHistory(c *gin.Context) {
	siteID := c.Query("site_id")
	if siteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少場站ID參數"})
		return
	}

	if !config.IsValidSite(siteID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的場站ID"})
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	limitStr := c.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的limit參數"})
		return
	}

	// 解析日期
	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無效的開始日期格式"})
			return
		}
	} else {
		startDate = time.Now().AddDate(0, 0, -30) // 預設30天前
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無效的結束日期格式"})
			return
		}
	} else {
		endDate = time.Now()
	}

	dataList, err := h.LoadModel.GetHistory(siteID, startDate, endDate, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"site_id":    siteID,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"count":      len(dataList),
		"data":       dataList,
	})
}

// GetSummary 獲取彙總統計
func (h *Handler) GetSummary(c *gin.Context) {
	// 獲取所有場站最新數據
	solarData, err := h.SolarModel.GetAllLatest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	loadData, err := h.LoadModel.GetAllLatest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 計算總和
	var totalGeneration, totalLoad, totalCO2 float64
	for _, solar := range solarData {
		totalGeneration += solar.DailyGeneration
		totalCO2 += solar.CO2Reduction
	}

	for _, load := range loadData {
		totalLoad += load.LoadValue
	}

	c.JSON(http.StatusOK, gin.H{
		"total_generation":  totalGeneration,
		"total_load":        totalLoad,
		"total_co2_reduced": totalCO2,
		"site_count":        len(solarData),
		"timestamp":         time.Now(),
	})
}
