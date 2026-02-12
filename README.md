# 牟牛 (MouNiu) 股票分析系统

这是一个专业的股票数据抓取与技术指标分析系统，专为金融数据分析和量化策略研究设计。系统能够自动抓取股票 K 线数据、计算多种核心技术指标，并提供标准化的 API 接口。

## 🌸 项目简介

“牟牛”系统由大姐为你精心打造，旨在提供稳定、高效、专业的股票分析基础服务。系统采用时序数据库 QuestDB 存储海量价格数据，并通过定时任务自动化完成复杂的指标计算逻辑。

## 🚀 核心功能

- **实时/历史数据抓取**：集成新浪财经 API，支持 A 股、港股等多种标的的 K 线数据采集。
- **自动化指标计算**：
  - **MACD (12, 26, 9)**：趋势追踪利器。
  - **布林带 (Bollinger Bands, 20, 2)**：波动率分析。
  - **神奇九转 (TD9)**：发现超买超跌的转折点。
  - **RSI (14)**：相对强弱指标。
  - **SMA (20)**：简单移动平均线。
- **公告爬取**：自动监控并存储上市公司的最新公告信息。
- **专业 API 服务**：基于 Gin 框架，提供完善的 RESTful 接口及 Swagger 文档支持。

## 🛠 技术栈

- **语言**：Go (Golang) 1.25+
- **框架**：Gin Web Framework
- **数据库**：QuestDB (高性能时序数据库)
- **ORM**：GORM
- **文档**：Swagger (swag)
- **调度**：robfig/cron/v3

## 📦 快速开始

### 1. 环境准备
确保你已经安装了 Go 语言环境和 QuestDB 数据库。

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 生成 API 文档
```bash
swag init -g main.go --parseDependency --parseInternal
```

### 4. 运行系统
```bash
go run main.go
```

## 📊 数据库查询指南 (QuestDB)

系统推荐使用 QuestDB 的 `LATEST ON` 语法来获取股票的最新状态：

```sql
-- 获取特定股票的最新的 K 线数据
SELECT * FROM candle_stick_data 
WHERE 股票代码 = 'SH600519' 
LATEST ON timestamp PARTITION BY 股票代码;

-- 获取所有股票的最新的技术指标分析
SELECT * FROM stock_indicators 
LATEST ON timestamp PARTITION BY stock_code;
```

## 🕒 定时任务配置

- **每 5 分钟**：自动抓取最新的股票价格数据。
- **每 10 分钟**：全量重新计算所有股票的技术指标。
- **每 3 小时**：抓取最新的上市公司公告。

## 📖 API 文档

系统启动后，可以通过以下地址访问交互式 Swagger 文档：
`http://localhost:8080/swagger/index.html`


