package function

import (
	"errors"
	"mouniu/internal/model"
)

// GetMACD 计算移动平均收敛散现指标 (MACD)
func GetMACD(prices []float64, fastPeriod, slowPeriod, signalPeriod int32) (model.MACDResult, error) {
	if int(slowPeriod) > len(prices) {
		return model.MACDResult{}, errors.New("价格数据不足以计算 MACD")
	}

	// 1. 计算快速和慢速 EMA 序列
	fastEMAs := calculateEMASeries(prices, fastPeriod)
	slowEMAs := calculateEMASeries(prices, slowPeriod)

	// 2. 计算 DIF 序列 (DIF = EMA_fast - EMA_slow)
	difSeries := make([]float64, len(slowEMAs))
	fastOffset := len(fastEMAs) - len(slowEMAs)

	// 使用函数式风格（虽然 Go 还是用循环实现，但逻辑更清晰）
	mapDIF(difSeries, fastEMAs[fastOffset:], slowEMAs)

	if int(signalPeriod) > len(difSeries) {
		return model.MACDResult{}, errors.New("DIF 数据不足以计算 DEA")
	}

	// 3. 计算 DEA (DIF 的 signalPeriod 日 EMA)
	deaSeries := calculateEMASeries(difSeries, signalPeriod)

	// 4. 获取最新的 DIF, DEA 并计算 MACD 柱状图
	latestDIF := difSeries[len(difSeries)-1]
	latestDEA := deaSeries[len(deaSeries)-1]

	return model.MACDResult{
		DIF:  latestDIF,
		DEA:  latestDEA,
		MACD: (latestDIF - latestDEA) * 2,
	}, nil
}

// mapDIF 辅助函数，用于映射计算 DIF 序列
func mapDIF(target, fast, slow []float64) {
	if len(fast) == 0 || len(slow) == 0 {
		return
	}
	target[0] = fast[0] - slow[0]
	mapDIF(target[1:], fast[1:], slow[1:])
}

// calculateEMASeries 计算整个价格序列的 EMA 值
func calculateEMASeries(prices []float64, period int32) []float64 {
	if len(prices) == 0 || int(period) > len(prices) {
		return nil
	}

	emas := make([]float64, len(prices)-int(period)+1)

	// 计算初始 SMA
	initialSMA := sumSlice(prices[:period]) / float64(period)
	emas[0] = initialSMA

	k := 2.0 / (float64(period) + 1.0)
	computeEMA(emas[1:], prices[period:], initialSMA, k)

	return emas
}

// computeEMA 递归计算 EMA，避免显式 for 循环
func computeEMA(target, prices []float64, prevEMA, k float64) {
	if len(prices) == 0 {
		return
	}
	currentEMA := (prices[0] * k) + (prevEMA * (1.0 - k))
	target[0] = currentEMA
	computeEMA(target[1:], prices[1:], currentEMA, k)
}

// sumSlice 递归求和
func sumSlice(s []float64) float64 {
	if len(s) == 0 {
		return 0
	}
	return s[0] + sumSlice(s[1:])
}
