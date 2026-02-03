package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetLatestReserve 獲取最新一天的備轉資料
func (h *Handler) GetLatestReserve(c *gin.Context) {
	dataList, err := h.TaipowerModel.GetLatest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(dataList) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到備轉資料"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"date":  dataList[0].TranDate.Format("2006-01-02"),
		"count": len(dataList),
		"data":  dataList,
	})
}

// GetReserveByDate 獲取特定日期的備轉資料
func (h *Handler) GetReserveByDate(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少日期參數"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的日期格式，請使用 YYYY-MM-DD"})
		return
	}

	dataList, err := h.TaipowerModel.GetByDate(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(dataList) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到該日期的備轉資料"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"date":  date.Format("2006-01-02"),
		"count": len(dataList),
		"data":  dataList,
	})
}

// GetReserveHistory 獲取歷史備轉資料
func (h *Handler) GetReserveHistory(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	limitStr := c.DefaultQuery("limit", "1000")

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

	dataList, err := h.TaipowerModel.GetHistory(startDate, endDate, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"count":      len(dataList),
		"data":       dataList,
	})
}

// GetReserveStatistics 獲取統計資訊
func (h *Handler) GetReserveStatistics(c *gin.Context) {
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的日期格式"})
		return
	}

	stats, err := h.TaipowerModel.GetStatistics(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"date":       date.Format("2006-01-02"),
		"statistics": stats,
	})
}

// GetReserveByHour 獲取特定時段的備轉資料
func (h *Handler) GetReserveByHour(c *gin.Context) {
	dateStr := c.Query("date")
	hourStr := c.Query("hour")

	if dateStr == "" || hourStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少日期或小時參數"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的日期格式"})
		return
	}

	hour, err := strconv.Atoi(hourStr)
	if err != nil || hour < 0 || hour > 23 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的小時參數（0-23）"})
		return
	}

	data, err := h.TaipowerModel.GetByHour(date, hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if data == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到該時段的備轉資料"})
		return
	}

	c.JSON(http.StatusOK, data)
}
