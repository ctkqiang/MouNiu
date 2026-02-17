# 牟牛 (MouNiu) 股票分析系统

牟牛 (MouNiu) 是一个基于 Go 语言开发的高性能股票数据抓取与技术指标分析系统。该系统集成了实时行情采集、多维技术指标计算、自动化公告监控以及标准化的 RESTful API 服务，旨在为量化交易和金融数据研究提供坚实的数据基础。

## 项目架构

系统采用分层架构设计，确保了代码的可维护性与扩展性：

1.  **数据采集层 (Data Acquisition)**：通过高性能并发机制从新浪财经等外部 API 抓取实时 K 线数据及上市公司公告。
2.  **逻辑处理层 (Business Logic)**：基于递归算法和高效切片操作实现技术指标计算，确保逻辑严谨且运行高效。
3.  **持久化层 (Persistence)**：使用高性能时序数据库 QuestDB 存储海量价格数据和指标结果，通过 LATEST ON 等特性优化查询性能。
4.  **接口层 (API Service)**：基于 Gin 框架提供高性能的 RESTful 接口，并集成 Swagger 生成标准化文档。

## 目录结构说明

```text
MouNiu/
├── contract/             # 协议定义文件 (如 Protobuf)
├── docs/                 # 自动生成的 Swagger API 文档
├── internal/
│   ├── config/           # 全局配置、API 终端定义及股票代码列表
│   ├── crons/            # 定时任务调度逻辑 (价格更新、指标计算、公告采集)
│   ├── database/         # 数据库连接驱动 (QuestDB, MySQL)
│   ├── function/         # 核心技术指标算法实现 (MACD, Bollinger, TD9, RSI, SMA)
│   ├── model/            # 领域模型与数据结构定义
│   ├── routes/           # API 路由注册与控制器 Handler 实现
│   ├── services/         # 业务逻辑服务 (数据抓取、指标持久化逻辑)
│   └── utilities/        # 通用工具库 (格式化、日志管理)
├── test/                 # 单元测试与集成测试
├── main.go               # 应用程序入口
└── README.md             # 项目说明文档
```

## 技术指标实现详情

系统实现了多种核心金融技术指标，计算逻辑均经过严格测试：

*   **MACD (12, 26, 9)**：采用递归指数移动平均 (EMA) 算法计算 DIF 和 DEA，并生成 MACD 柱状图数据。
*   **布林带 (Bollinger Bands, 20, 2)**：计算 20 日简单移动平均线作为中轨，结合标准差计算上下轨，用于分析股价波动区间。
*   **神奇九转 (TD9)**：通过对比连续 4 日的收盘价，识别趋势中的转折点，支持买入结构和卖出结构的自动判定。
*   **RSI (14)**：相对强弱指标，用于衡量股价变动的速度和变化，判定超买或超跌状态。
*   **SMA (20)**：简单移动平均线，反映股价的长期运行趋势。

## 数据同步与调度机制

系统内置了完善的 Cron 调度任务，确保数据的实时性与一致性：

*   **行情更新 (5分钟/次)**：并发抓取 `symbols.txt` 中配置的所有股票实时行情。
*   **指标重算 (10分钟/次)**：对历史价格序列进行全量分析，将计算结果存入 `stock_indicators` 表。
*   **公告监控 (3小时/次)**：自动采集并同步深交所、港交所等市场的上市公司公告。

## 数据库查询优化 (QuestDB)

针对时序数据的特性，系统在查询层进行了针对性优化。推荐使用 `LATEST ON` 语法获取最新状态，以减少全表扫描：

```sql
-- 获取特定股票最新的行情记录
SELECT * FROM candle_stick_data 
WHERE 股票代码 = 'SH600519' 
LATEST ON timestamp PARTITION BY 股票代码;

-- 获取所有股票最新的指标分析结果
SELECT * FROM stock_indicators 
LATEST ON timestamp PARTITION BY stock_code;
```

## API 文档与接入

系统集成 Swagger (OpenAPI 2.0)，启动后可通过以下地址访问交互式文档：

`http://localhost:8080/swagger/index.html`

接口涵盖了股票历史数据查询、实时价格获取、最新指标分析结果以及公告信息检索。

### **常用 API 调用示例 (cURL)**

以下是一些常用的接口调用示例，你可以直接在终端中运行：

*   **获取所有股票的最新的技术指标分析**：
    ```bash
    curl http://localhost:8080/api/analysis
    ```
*   **获取特定股票 (如 贵州茅台 SH600519) 的最新指标**：
    ```bash
    curl http://localhost:8080/api/analysis/SH600519
    ```
*   **获取特定股票的所有历史 K 线数据**：
    ```bash
    curl http://localhost:8080/api/SH600519
    ```
*   **获取特定股票的当前最新价格记录**：
    ```bash
    curl http://localhost:8080/api/current/SH600519
    ```
*   **分页获取所有上市公司的公告**：
    ```bash
    curl "http://localhost:8080/api/announcement/all?page=1&pageSize=10"
    ```
*   **查询特定股票代码相关的公告**：
    ```bash
    curl "http://localhost:8080/api/announcement/600519?page=1&pageSize=5"
    ```

## 开发与部署

### 环境依赖
*   Go 1.25 或更高版本
*   QuestDB 8.0+

### 编译运行
1.  安装依赖：`go mod tidy`
2.  更新文档：`swag init -g main.go --parseDependency --parseInternal`
3.  启动服务：`go run main.go`

---

### 🌐 全球捐赠通道

#### 国内用户

<div align="center" style="margin: 40px 0">

<div align="center">
<table>
<tr>
<td align="center" width="300">
<img src="https://github.com/ctkqiang/ctkqiang/blob/main/assets/IMG_9863.jpg?raw=true" width="200" />
<br />
<strong>🔵 支付宝</strong>（小企鹅在收金币哟~）
</td>
<td align="center" width="300">
<img src="https://github.com/ctkqiang/ctkqiang/blob/main/assets/IMG_9859.JPG?raw=true" width="200" />
<br />
<strong>🟢 微信支付</strong>（小绿龙在收金币哟~）
</td>
</tr>
</table>
</div>
</div>

#### 国际用户

<div align="center" style="margin: 40px 0">
  <a href="https://qr.alipay.com/fkx19369scgxdrkv8mxso92" target="_blank">
    <img src="https://img.shields.io/badge/Alipay-全球支付-00A1E9?style=flat-square&logo=alipay&logoColor=white&labelColor=008CD7">
  </a>
  
  <a href="https://ko-fi.com/F1F5VCZJU" target="_blank">
    <img src="https://img.shields.io/badge/Ko--fi-买杯咖啡-FF5E5B?style=flat-square&logo=ko-fi&logoColor=white">
  </a>
  
  <a href="https://www.paypal.com/paypalme/ctkqiang" target="_blank">
    <img src="https://img.shields.io/badge/PayPal-安全支付-00457C?style=flat-square&logo=paypal&logoColor=white">
  </a>
  
  <a href="https://donate.stripe.com/00gg2nefu6TK1LqeUY" target="_blank">
    <img src="https://img.shields.io/badge/Stripe-企业级支付-626CD9?style=flat-square&logo=stripe&logoColor=white">
  </a>
</div>

---

### 📌 开发者社交图谱

#### 技术交流

<div align="center" style="margin: 20px 0">
  <a href="https://github.com/ctkqiang" target="_blank">
    <img src="https://img.shields.io/badge/GitHub-开源仓库-181717?style=for-the-badge&logo=github">
  </a>
  
  <a href="https://stackoverflow.com/users/10758321/%e9%92%9f%e6%99%ba%e5%bc%ba" target="_blank">
    <img src="https://img.shields.io/badge/Stack_Overflow-技术问答-F58025?style=for-the-badge&logo=stackoverflow">
  </a>
  
  <a href="https://www.linkedin.com/in/ctkqiang/" target="_blank">
    <img src="https://img.shields.io/badge/LinkedIn-职业网络-0A66C2?style=for-the-badge&logo=linkedin">
  </a>
</div>

#### 社交互动

<div align="center" style="margin: 20px 0">
  <a href="https://www.instagram.com/ctkqiang" target="_blank">
    <img src="https://img.shields.io/badge/Instagram-生活瞬间-E4405F?style=for-the-badge&logo=instagram">
  </a>
  
  <a href="https://twitch.tv/ctkqiang" target="_blank">
    <img src="https://img.shields.io/badge/Twitch-技术直播-9146FF?style=for-the-badge&logo=twitch">
  </a>
  
  <a href="https://github.com/ctkqiang/ctkqiang/blob/main/assets/IMG_9245.JPG?raw=true" target="_blank">
    <img src="https://img.shields.io/badge/微信公众号-钟智强-07C160?style=for-the-badge&logo=wechat">
  </a>
</div>