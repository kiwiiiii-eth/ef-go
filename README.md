# VPP (虛擬電廠) API - Go版本

這是一個使用 Go 語言重寫的虛擬電廠管理系統 API，專為 Zeabur 平台部署而優化。

## 功能特點

- ⚡ **高性能**: 使用 Go 語言和 Gin 框架，提供更快的響應速度
- 🏗️ **清晰架構**: 採用分層架構設計（handlers, models, database）
- 📊 **完整功能**:
  - VPP 數據管理（太陽能、負載、儲能）
  - 台電備轉資料查詢
  - 樹莓派數據上傳
  - 自動化數據收集
- 🔌 **PostgreSQL**: 使用 PostgreSQL 資料庫
- 🚀 **Zeabur 部署**: 一鍵部署到 Zeabur 平台
- 🕐 **自動收集**: 定時收集太陽能和台電備轉資料

## 專案結構

```
vpp-go/
├── cmd/
│   └── api/
│       └── main.go              # 主程式入口
├── internal/
│   ├── config/
│   │   └── config.go            # 配置管理
│   ├── database/
│   │   └── database.go          # 資料庫連接
│   ├── models/
│   │   ├── solar.go             # 太陽能數據模型
│   │   ├── load.go              # 負載數據模型
│   │   └── taipower.go          # 台電備轉資料模型
│   ├── handlers/
│   │   ├── handler.go           # 處理器基礎
│   │   ├── vpp.go               # VPP API 處理器
│   │   ├── taipower.go          # 台電 API 處理器
│   │   └── upload.go            # 上傳 API 處理器
│   └── collectors/
│       ├── solar_collector.go   # 太陽能數據收集器
│       └── taipower_collector.go # 台電數據收集器
├── pkg/
│   └── utils/                   # 工具函數
├── go.mod                       # Go 模組定義
├── go.sum                       # Go 依賴鎖定
├── Dockerfile                   # Docker 構建文件
├── Makefile                     # 構建腳本
├── .env.example                 # 環境變數示例
├── zbpack.json                  # Zeabur 配置
└── README.md                    # 本文件
```

## 技術棧

- **語言**: Go 1.21+
- **Web 框架**: Gin
- **資料庫**: PostgreSQL
- **HTTP 客戶端**: net/http
- **定時任務**: time.Ticker
- **依賴管理**: Go Modules

## 快速開始

### 1. 本地開發

#### 前置要求

- Go 1.21 或更高版本
- PostgreSQL 資料庫
- Make (可選)

#### 安裝依賴

```bash
cd vpp-go
go mod download
```

#### 配置環境變數

複製 `.env.example` 為 `.env` 並填入配置：

```bash
cp .env.example .env
```

編輯 `.env` 文件：

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=vpp_db
PORT=8080
```

#### 初始化資料庫

確保資料庫表已創建（使用原 Flask 專案的 init_db.py 或手動創建）。

#### 運行應用

```bash
# 使用 go run
go run ./cmd/api/main.go

# 或使用 make
make run

# 或先構建再運行
make build
./bin/vpp-api
```

應用將在 `http://localhost:8080` 啟動。

### 2. Docker 部署

#### 構建 Docker 映像

```bash
make docker-build
# 或
docker build -t vpp-go:latest .
```

#### 運行 Docker 容器

```bash
make docker-run
# 或
docker run -p 8080:8080 --env-file .env vpp-go:latest
```

### 3. Zeabur 部署

#### 步驟 1: 準備資料庫

1. 在 Zeabur 創建 PostgreSQL 服務
2. 獲取資料庫連接信息

#### 步驟 2: 部署應用

1. 將代碼推送到 GitHub 倉庫
2. 在 Zeabur 控制台點擊 "New Project"
3. 選擇 "Deploy from GitHub"
4. 選擇你的倉庫和分支
5. Zeabur 會自動檢測到 Go 專案並開始構建

#### 步驟 3: 配置環境變數

在 Zeabur 控制台的 Environment Variables 頁面添加：

```
DB_HOST=your-db-host.zeabur.internal
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-db-password
DB_NAME=vpp_db
PORT=8080
GIN_MODE=release
```

#### 步驟 4: 連接資料庫

在 Zeabur 控制台將 PostgreSQL 服務連接到你的應用。

#### 步驟 5: 部署完成

應用部署完成後，你會獲得一個公開的 URL。

## API 端點

### 根路由

```
GET /
```

返回 API 信息和可用端點列表。

### VPP 路由

#### 即時數據

- `GET /api/vpp/realdata` - 獲取所有場站即時數據
- `GET /api/vpp/realdata/:site_id` - 獲取特定場站即時數據

#### 太陽能數據

- `GET /api/vpp/solar/latest` - 獲取最新太陽能數據
  - 參數: `site_id` (可選)
- `GET /api/vpp/solar/history` - 獲取歷史太陽能數據
  - 參數: `site_id` (必須), `start_date`, `end_date`, `limit`

#### 負載數據

- `GET /api/vpp/load/latest` - 獲取最新負載數據
  - 參數: `site_id` (可選)
- `GET /api/vpp/load/history` - 獲取歷史負載數據
  - 參數: `site_id` (必須), `start_date`, `end_date`, `limit`

#### 統計彙總

- `GET /api/vpp/summary` - 獲取彙總統計

### 台電備轉資料路由

- `GET /api/taipower/reserve/latest` - 獲取最新一天備轉資料
- `GET /api/taipower/reserve/date` - 獲取特定日期備轉資料
  - 參數: `date` (YYYY-MM-DD)
- `GET /api/taipower/reserve/history` - 獲取歷史備轉資料
  - 參數: `start_date`, `end_date`, `limit`
- `GET /api/taipower/reserve/statistics` - 獲取統計資訊
  - 參數: `date` (可選)
- `GET /api/taipower/reserve/hour` - 獲取特定時段備轉資料
  - 參數: `date` (YYYY-MM-DD), `hour` (0-23)

### 上傳路由

- `POST /api/upload` - 樹莓派數據上傳

## 數據收集器

### 太陽能數據收集器

自動每 15 分鐘從義鴻太陽能 API 收集數據。

配置環境變數：
```
YIHONG_API_URL=https://api.yihong-solar.com/data
YIHONG_USERNAME=your-username
YIHONG_PASSWORD=your-password
```

### 台電備轉資料收集器

自動每天凌晨 2 點收集前一天的台電備轉資料。

配置環境變數：
```
TAIPOWER_URL=https://www.taipower.com.tw
```

## 場站 ID

系統支援三個場站：

- `north` - 北部場站
- `central` - 中部場站
- `south` - 南部場站

## 開發

### 運行測試

```bash
make test
# 或
go test -v ./...
```

### 代碼格式化

```bash
make format
# 或
go fmt ./...
```

### 代碼檢查

```bash
make lint
# 或
golangci-lint run
```

## 性能優化

相比 Flask 版本，Go 版本提供以下優勢：

1. **更快的響應速度**: Go 的編譯型特性和高效的 goroutine
2. **更低的內存佔用**: Go 的垃圾回收機制更高效
3. **更好的並發處理**: 原生的 goroutine 支持
4. **更快的啟動時間**: 編譯後的二進制文件啟動速度快
5. **更小的部署包**: 單一二進制文件，無需依賴

## 與 Flask 版本的差異

1. **語言**: Python → Go
2. **框架**: Flask → Gin
3. **性能**: 更快的響應速度和更低的資源佔用
4. **部署**: 單一二進制文件，部署更簡單
5. **並發**: 使用 goroutine 處理並發請求

## 故障排除

### 無法連接資料庫

檢查環境變數配置是否正確：
```bash
echo $DB_HOST
echo $DB_PORT
```

### 端口已被佔用

修改 `.env` 文件中的 `PORT` 變數。

### 數據收集器無法運行

檢查 API URL 和認證信息是否正確配置。

## 貢獻

歡迎提交 Issue 和 Pull Request！

## 授權

[添加你的授權信息]

## 聯繫方式

[添加你的聯繫方式]
