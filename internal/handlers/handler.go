package handlers

import (
	"database/sql"
	"vpp-go/internal/models"
)

// Handler 處理器結構
type Handler struct {
	DB            *sql.DB
	SolarModel    *models.SolarDataModel
	LoadModel     *models.LoadDataModel
	TaipowerModel *models.TaipowerReserveModel
}

// NewHandler 創建新的處理器
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		DB:            db,
		SolarModel:    models.NewSolarDataModel(db),
		LoadModel:     models.NewLoadDataModel(db),
		TaipowerModel: models.NewTaipowerReserveModel(db),
	}
}
