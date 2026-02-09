package function

import (
	"errors"
	"math"
)

// GetExponentialMovingAverage 计算指数移动平均线（EMA）
// EMA是一种技术分析指标，给予近期价格更高的权重，反应更灵敏
// 公式：EMA_today = (Price_today × k) + (EMA_yesterday × (1 - k))
// 其中 k = 2 / (period + 1) 为平滑系数
//
// 参数：
//   - lastClosingPrice: 最新收盘价，必须为正数
//   - previousEMA: 前一周期EMA值，若无历史数据可使用收盘价作为初始值
//   - period: EMA计算周期，必须大于0，常用值：5(周)、10、20、30、50、100、200
//
// 返回：
//   - float64: 计算得到的EMA值
//   - error: 如果参数无效返回错误
//
// 示例：
//
//	假设计算5日EMA，今日收盘价100，昨日EMA 95
//	k = 2/(5+1) = 0.333
//	EMA = (100×0.333) + (95×0.667) = 96.665
func GetExponentialMovingAverage(lastClosingPrice float64, previousEMA float64, period int32) (float64, error) {
	if period <= 0 {
		return 0, errors.New("EMA计算周期必须大于0")
	}

	if lastClosingPrice <= 0 {
		return 0, errors.New("收盘价必须为正数")
	}

	if previousEMA < 0 {
		return 0, errors.New("前一EMA值不能为负数")
	}

	if math.IsNaN(lastClosingPrice) || math.IsInf(lastClosingPrice, 0) {
		return 0, errors.New("收盘价不是有效数值")
	}

	if math.IsNaN(previousEMA) || math.IsInf(previousEMA, 0) {
		return 0, errors.New("前一EMA值不是有效数值")
	}

	// 计算平滑系数 k
	k := 2.0 / (float64(period) + 1.0)

	// 返回当前EMA值：今日收盘价权重 k，前一EMA 权重 (1-k)
	return (lastClosingPrice * k) + (previousEMA * (1.0 - k)), nil
}

// GetSimpleMovingAverage 计算简单移动平均线（SMA）
// SMA是技术分析中最基础的指标，取最近N个周期的算术平均值
// 公式：SMA = (P1 + P2 + ... + Pn) / n
// 其中 P1 到 Pn 为最近n个周期的价格
//
// 参数：
//   - prices: 价格序列切片，按时间顺序排列（最新价格在最后）
//   - period: SMA计算周期，必须大于0且小于等于价格序列长度
//
// 返回：
//   - float64: 计算得到的SMA值
//   - error: 如果参数无效返回错误
//
// 示例：
//
//	假设有最近5日收盘价：[95, 96, 97, 98, 100]
//	SMA = (95+96+97+98+100)/5 = 97.2
func GetSimpleMovingAverage(prices []float64, period int32) (float64, error) {
	if period <= 0 {
		return 0, errors.New("SMA计算周期必须大于0")
	}

	if len(prices) == 0 {
		return 0, errors.New("价格序列不能为空")
	}

	if int(period) > len(prices) {
		return 0, errors.New("计算周期不能大于价格序列长度")
	}

	for _, price := range prices {
		if price <= 0 {
			return 0, errors.New("价格序列中的值必须为正数")
		}

		if math.IsNaN(price) || math.IsInf(price, 0) {
			return 0, errors.New("价格序列包含无效数值")
		}
	}

	startIndex := len(prices) - int(period)
	sum := 0.0

	for i := startIndex; i < len(prices); i++ {
		sum += prices[i]
	}

	return sum / float64(period), nil
}
