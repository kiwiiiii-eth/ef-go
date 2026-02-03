package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadRequest 上傳請求結構
type UploadRequest struct {
	SiteID    string  `json:"site_id" binding:"required"`
	Timestamp string  `json:"timestamp"`
	Data      RawData `json:"data" binding:"required"`
}

// RawData 原始數據結構（保留給樹莓派使用）
type RawData struct {
	Value interface{} `json:"value"`
}

// UploadData 處理樹莓派數據上傳
func (h *Handler) UploadData(c *gin.Context) {
	var req UploadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求數據"})
		return
	}

	// 解析時間戳
	var timestamp time.Time
	var err error
	if req.Timestamp != "" {
		timestamp, err = time.Parse(time.RFC3339, req.Timestamp)
		if err != nil {
			timestamp = time.Now()
		}
	} else {
		timestamp = time.Now()
	}

	// 這裡可以根據需要將數據插入到 stu 表或其他表
	// 目前保留原有的上傳功能，可以根據實際需求修改

	// 假設插入到一個通用的上傳表
	query := `
		INSERT INTO stu (site_id, timestamp, data)
		VALUES ($1, $2, $3)
	`

	_, err = h.DB.Exec(query, req.SiteID, timestamp, req.Data.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "數據保存失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "數據上傳成功",
		"site_id":   req.SiteID,
		"timestamp": timestamp.Format(time.RFC3339),
	})
}
