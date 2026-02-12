package function

import (
	"errors"
)

// GetExponentialMovingAverage 计算指数移动平均线（EMA）
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

	k := 2.0 / (float64(period) + 1.0)
	return (lastClosingPrice * k) + (previousEMA * (1.0 - k)), nil
}

// GetSimpleMovingAverage 计算简单移动平均线（SMA）
func GetSimpleMovingAverage(prices []float64, period int32) (float64, error) {
	if period <= 0 || len(prices) == 0 || int(period) > len(prices) {
		return 0, errors.New("无效的 SMA 计算参数")
	}

	startIndex := len(prices) - int(period)
	sum := recursiveSum(prices[startIndex:])
	return sum / float64(period), nil
}

// recursiveSum 递归求和，避免 for 循环
func recursiveSum(s []float64) float64 {
	if len(s) == 0 {
		return 0
	}
	return s[0] + recursiveSum(s[1:])
}
