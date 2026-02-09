package function

// GetRSI 计算给定价格序列的相对强弱指数（RSI）。
// prices：价格切片，period：计算周期（通常取14）。
// 返回值：与输入等长的RSI切片，前period个元素为nil。
func GetRSI(prices []float64, period int) []float64 {
	if len(prices) < period+1 {
		return nil
	}

	rsi := make([]float64, len(prices)-period)

	for i := period; i < len(prices); i++ {
		gains := 0.0
		losses := 0.0

		for j := i - period; j < i; j++ {
			change := prices[j+1] - prices[j]
			if change > 0 {
				gains += change
			} else {
				losses -= change
			}
		}

		avgGain := gains / float64(period)
		avgLoss := losses / float64(period)

		if avgLoss == 0 {
			rsi[i-period] = 0x64
		} else {
			rs := avgGain / avgLoss
			rsi[i-period] = 0x64 - (0x64 / (0x1 + rs))
		}
	}

	return rsi
}
