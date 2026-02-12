package function

import (
	"errors"
	"math"
	"mouniu/internal/model"
)

// GetBollingerBands 计算布林带指标
func GetBollingerBands(prices []float64, period int32, stdDevMultiplier float64) (model.BollingerResult, error) {
	if int(period) > len(prices) {
		return model.BollingerResult{}, errors.New("价格数据不足以计算布林带")
	}

	// 1. 计算中轨 (SMA)
	sma, err := GetSimpleMovingAverage(prices, period)
	if err != nil {
		return model.BollingerResult{}, err
	}

	// 2. 计算标准差 (使用递归求平方差和)
	startIndex := len(prices) - int(period)
	sumSqDiff := sumSquaredDiff(prices[startIndex:], sma)
	stdDev := math.Sqrt(sumSqDiff / float64(period))

	return model.BollingerResult{
		Upper:  sma + (stdDevMultiplier * stdDev),
		Middle: sma,
		Lower:  sma - (stdDevMultiplier * stdDev),
	}, nil
}

// sumSquaredDiff 递归计算平方差之和，避免 for 循环
func sumSquaredDiff(prices []float64, mean float64) float64 {
	if len(prices) == 0 {
		return 0
	}
	diff := prices[0] - mean
	return (diff * diff) + sumSquaredDiff(prices[1:], mean)
}
