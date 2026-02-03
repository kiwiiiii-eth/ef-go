# Zeabur 部署指南

本指南將詳細介紹如何將 VPP API Go 版本部署到 Zeabur 平台。

## 前置準備

1. GitHub 帳號
2. Zeabur 帳號（使用 GitHub 登入）
3. 本專案的 Git 倉庫

## 部署步驟

### 第一步：準備代碼倉庫

#### 1. 初始化 Git 倉庫（如果尚未初始化）

```bash
cd vpp-go
git init
git add .
git commit -m "Initial commit: VPP API Go version"
```

#### 2. 推送到 GitHub

```bash
# 創建 GitHub 倉庫後
git remote add origin https://github.com/your-username/vpp-go.git
git branch -M main
git push -u origin main
```

### 第二步：創建 PostgreSQL 資料庫

#### 1. 登入 Zeabur

前往 [Zeabur 控制台](https://dash.zeabur.com)

#### 2. 創建新專案

- 點擊 "New Project"
- 輸入專案名稱（例如：vpp-system）

#### 3. 添加 PostgreSQL 服務

- 在專案頁面點擊 "Add Service"
- 選擇 "Prebuilt" → "PostgreSQL"
- 等待 PostgreSQL 服務部署完成

#### 4. 獲取資料庫連接信息

- 點擊 PostgreSQL 服務卡片
- 進入 "Connect" 標籤
- 記錄以下信息：
  - Host
  - Port
  - Username
  - Password
  - Database Name

#### 5. 初始化資料庫表

使用以下方式之一初始化資料庫表：

**方式 1: 使用原 Flask 專案的腳本**

```bash
# 在原 Flask 專案目錄
export DATABASE_URL="postgresql://user:password@host:port/dbname"
python scripts/init_db.py
```

**方式 2: 手動執行 SQL**

連接到 Zeabur 的 PostgreSQL 並執行創建表的 SQL 語句。

### 第三步：部署 Go 應用

#### 1. 在 Zeabur 添加 Git 服務

- 在同一專案中點擊 "Add Service"
- 選擇 "Git"
- 授權 GitHub 訪問
- 選擇你的倉庫：`your-username/vpp-go`
- 選擇分支：`main`

#### 2. 配置構建設置（Zeabur 會自動檢測）

Zeabur 會自動檢測到 `zbpack.json` 文件並使用配置：

```json
{
  "build_command": "go build -o main ./cmd/api",
  "start_command": "./main",
  "install_command": "go mod download"
}
```

如果需要手動配置，可以在服務設置中修改。

#### 3. 配置環境變數

在服務設置的 "Environment Variables" 標籤添加以下變數：

**基本配置**:
```
PORT=8080
GIN_MODE=release
```

**資料庫配置**（使用 PostgreSQL 服務的連接信息）:
```
DB_HOST=postgresql.zeabur.internal
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=<your-database-password>
DB_NAME=vpp_db
DB_SSLMODE=disable
```

**外部 API 配置**（可選，如果需要自動收集數據）:
```
YIHONG_API_URL=https://api.yihong-solar.com/data
YIHONG_USERNAME=<your-username>
YIHONG_PASSWORD=<your-password>
TAIPOWER_URL=https://www.taipower.com.tw
```

**場站配置**:
```
SITE_NORTH=north
SITE_CENTRAL=central
SITE_SOUTH=south
```

#### 4. 連接服務

在 Zeabur 控制台中：
- 點擊 Go 應用服務
- 進入 "Connect" 標籤
- 點擊 "Connect to PostgreSQL"
- 選擇你之前創建的 PostgreSQL 服務

這會自動設置內部網絡連接。

#### 5. 開始部署

- 點擊 "Deploy" 按鈕
- Zeabur 會自動：
  1. 克隆代碼
  2. 安裝 Go 依賴
  3. 構建應用
  4. 啟動服務

#### 6. 檢查部署日誌

在 "Logs" 標籤查看部署日誌，確保沒有錯誤。

### 第四步：配置域名

#### 1. 獲取自動生成的域名

部署完成後，Zeabur 會自動生成一個域名，例如：
```
https://vpp-go-xxx.zeabur.app
```

#### 2. 配置自定義域名（可選）

如果你有自己的域名：

1. 進入服務的 "Domains" 標籤
2. 點擊 "Add Domain"
3. 輸入你的域名（例如：api.yourdomain.com）
4. 在你的 DNS 提供商添加 CNAME 記錄：
   ```
   api.yourdomain.com → vpp-go-xxx.zeabur.app
   ```
5. 等待 DNS 傳播（通常幾分鐘）
6. Zeabur 會自動配置 SSL 證書

### 第五步：測試部署

#### 1. 測試根端點

```bash
curl https://your-app.zeabur.app/
```

應該返回：
```json
{
  "message": "VPP (虛擬電廠) API - Go版本",
  "version": "1.0.0",
  "endpoints": {
    "upload": "/api/upload",
    "vpp": "/api/vpp/*",
    "taipower": "/api/taipower/*"
  }
}
```

#### 2. 測試 VPP API

```bash
# 獲取所有場站即時數據
curl https://your-app.zeabur.app/api/vpp/realdata

# 獲取特定場站數據
curl https://your-app.zeabur.app/api/vpp/realdata/north

# 獲取太陽能最新數據
curl https://your-app.zeabur.app/api/vpp/solar/latest
```

#### 3. 測試台電 API

```bash
# 獲取最新備轉資料
curl https://your-app.zeabur.app/api/taipower/reserve/latest

# 獲取特定日期
curl "https://your-app.zeabur.app/api/taipower/reserve/date?date=2024-01-01"
```

## 持續部署（CI/CD）

Zeabur 支持自動部署：

1. 每次推送到 GitHub 的 `main` 分支時，Zeabur 會自動重新部署
2. 可以在 Zeabur 控制台的 "Settings" 中配置其他分支

### 設置自動部署

```bash
# 開發新功能
git checkout -b feature/new-feature
# ... 進行修改 ...
git add .
git commit -m "Add new feature"
git push origin feature/new-feature

# 合併到 main 分支觸發自動部署
git checkout main
git merge feature/new-feature
git push origin main
```

## 監控和日誌

### 查看即時日誌

在 Zeabur 控制台：
1. 進入你的服務
2. 點擊 "Logs" 標籤
3. 可以看到即時的應用日誌

### 查看性能指標

在 "Metrics" 標籤可以查看：
- CPU 使用率
- 內存使用量
- 網絡流量
- 請求數量

## 擴展和優化

### 1. 垂直擴展

在 Zeabur 控制台調整資源配置：
- 進入 "Settings" → "Resources"
- 選擇更大的計劃

### 2. 環境變數管理

使用 Zeabur 的環境變數功能：
- 開發環境和生產環境使用不同的配置
- 敏感信息不要硬編碼在代碼中

### 3. 定時任務

如果需要運行定時任務（數據收集器），可以：

**選項 1**: 在主應用中啟動（已實現）
- 數據收集器會在應用啟動時自動運行

**選項 2**: 使用 Zeabur Cron Jobs（如果可用）
- 創建單獨的 Cron Job 服務
- 配置定時執行腳本

## 故障排除

### 問題 1: 無法連接資料庫

**解決方案**:
1. 檢查環境變數是否正確設置
2. 確保 PostgreSQL 服務和 Go 服務在同一個 Zeabur 專案中
3. 使用內部域名：`postgresql.zeabur.internal`

### 問題 2: 部署失敗

**解決方案**:
1. 檢查 "Logs" 標籤的錯誤信息
2. 確保 `go.mod` 和 `go.sum` 文件已提交
3. 確保 `zbpack.json` 配置正確

### 問題 3: 應用無法啟動

**解決方案**:
1. 檢查 PORT 環境變數是否設置
2. 確保應用監聽 `0.0.0.0:$PORT`
3. 檢查日誌中的錯誤信息

### 問題 4: 502 Bad Gateway

**解決方案**:
1. 檢查應用是否正常啟動
2. 確保應用監聽正確的端口
3. 查看日誌排查錯誤

## 成本估算

Zeabur 提供不同的定價方案：

- **Hobby 方案**: 免費額度，適合開發測試
- **Pro 方案**: 按使用量計費，適合生產環境

建議：
- 開發階段使用 Hobby 方案
- 生產環境使用 Pro 方案並配置適當的資源

## 安全建議

1. **環境變數**: 所有敏感信息都使用環境變數
2. **資料庫**: 使用強密碼，啟用 SSL 連接
3. **API 金鑰**: 定期輪換 API 金鑰
4. **HTTPS**: 使用 Zeabur 提供的自動 SSL 證書
5. **訪問控制**: 考慮添加 API 認證機制

## 備份策略

### 資料庫備份

Zeabur PostgreSQL 服務提供自動備份功能：
1. 進入 PostgreSQL 服務設置
2. 啟用自動備份
3. 定期下載備份到本地

### 手動備份

```bash
# 使用 pg_dump
pg_dump -h your-db-host.zeabur.app -U postgres -d vpp_db > backup.sql
```

## 升級和維護

### 升級 Go 版本

1. 更新 `go.mod` 中的 Go 版本
2. 更新 `Dockerfile` 中的 Go 版本
3. 推送到 GitHub 觸發重新部署

### 更新依賴

```bash
go get -u ./...
go mod tidy
git add go.mod go.sum
git commit -m "Update dependencies"
git push
```

## 支援

如有問題，請參考：
- [Zeabur 文檔](https://zeabur.com/docs)
- [Go Gin 文檔](https://gin-gonic.com/docs/)
- 專案 GitHub Issues

## 下一步

部署完成後，你可以：

1. 配置數據收集器自動運行
2. 設置監控和告警
3. 添加 API 認證
4. 配置 CDN 加速
5. 優化資料庫查詢性能

祝你部署順利！ 🚀
