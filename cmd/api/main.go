package main

import (
	"fmt"
	"log"
	"os"
	"vpp-go/internal/config"
	"vpp-go/internal/database"
	"vpp-go/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// 初始化配置
	cfg := config.Load()

	// 初始化資料庫連接
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("無法連接資料庫: %v", err)
	}
	defer db.Close()

	// 設置Gin模式
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 創建路由
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 創建處理器
	h := handlers.NewHandler(db)

	// 根路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "VPP (虛擬電廠) API - Go版本",
			"version": "1.0.0",
			"endpoints": gin.H{
				"upload":   "/api/upload",
				"vpp":      "/api/vpp/*",
				"taipower": "/api/taipower/*",
			},
		})
	})

	// API路由組
	api := r.Group("/api")
	{
		// 樹莓派數據上傳
		api.POST("/upload", h.UploadData)

		// VPP路由
		vpp := api.Group("/vpp")
		{
			// 即時數據
			vpp.GET("/realdata", h.GetAllRealtimeData)
			vpp.GET("/realdata/:site_id", h.GetSiteRealtimeData)

			// 太陽能數據
			vpp.GET("/solar/latest", h.GetLatestSolarData)
			vpp.GET("/solar/history", h.GetSolarHistory)

			// 負載數據
			vpp.GET("/load/latest", h.GetLatestLoadData)
			vpp.GET("/load/history", h.GetLoadHistory)

			// 統計彙總
			vpp.GET("/summary", h.GetSummary)
		}

		// 台電備轉資料路由
		taipower := api.Group("/taipower")
		{
			reserve := taipower.Group("/reserve")
			{
				reserve.GET("/latest", h.GetLatestReserve)
				reserve.GET("/date", h.GetReserveByDate)
				reserve.GET("/history", h.GetReserveHistory)
				reserve.GET("/statistics", h.GetReserveStatistics)
				reserve.GET("/hour", h.GetReserveByHour)
			}
		}
	}

	// 獲取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 啟動服務器
	log.Printf("VPP API 服務器啟動於端口 %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服務器啟動失敗: %v", err)
	}
}
