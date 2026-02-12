package function

// GetRSI 计算给定价格序列的相对强弱指数（RSI）。
func GetRSI(prices []float64, period int) []float64 {
	if len(prices) < period+1 {
		return nil
	}

	rsi := make([]float64, len(prices)-period)
	computeRSISeries(rsi, prices, period)
	return rsi
}

// computeRSISeries 递归计算 RSI 序列
func computeRSISeries(target []float64, prices []float64, period int) {
	if len(target) == 0 {
		return
	}

	// 计算当前位置的 RSI
	target[0] = calculateSingleRSI(prices[:period+1], period)

	// 递归处理下一个位置
	computeRSISeries(target[1:], prices[1:], period)
}

// calculateSingleRSI 计算单个点的 RSI 值
func calculateSingleRSI(p []float64, period int) float64 {
	gains, losses := computeGainsAndLosses(p, 0, 0)
	
	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	if avgLoss == 0 {
		return 100.0
	}
	
	rs := avgGain / avgLoss
	return 100.0 - (100.0 / (1.0 + rs))
}

// computeGainsAndLosses 递归计算涨跌幅之和
func computeGainsAndLosses(p []float64, gains, losses float64) (float64, float64) {
	if len(p) < 2 {
		return gains, losses
	}
	
	change := p[1] - p[0]
	if change > 0 {
		gains += change
	} else {
		losses -= change
	}
	
	return computeGainsAndLosses(p[1:], gains, losses)
}
