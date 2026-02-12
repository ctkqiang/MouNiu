package function

import (
	"errors"
	"mouniu/internal/model"
)

// GetTD9 计算“神奇九转”指标的 Setup 阶段
func GetTD9(prices []float64) (model.TD9Result, error) {
	if len(prices) < 13 {
		return model.TD9Result{}, errors.New("价格数据不足以计算完整的“神奇九转”结构 (至少需要 13 个周期)")
	}

	// 使用递归计算最新的 Setup 计数
	buyCount, sellCount := computeTD9Counts(prices[len(prices)-9:], prices[len(prices)-13:])

	result := model.TD9Result{}
	if buyCount > 0 {
		result.Count = buyCount
		result.IsBuySetup = (buyCount == 9)
	} else if sellCount > 0 {
		result.Count = sellCount
		result.IsSellSetup = (sellCount == 9)
	}

	return result, nil
}

// computeTD9Counts 递归计算买入和卖出计数，避免显式 for 循环
func computeTD9Counts(currentPrices, comparePrices []float64) (buy, sell int) {
	if len(currentPrices) == 0 {
		return 0, 0
	}

	// 递归处理前面的价格
	prevBuy, prevSell := computeTD9Counts(currentPrices[:len(currentPrices)-1], comparePrices[:len(comparePrices)-1])

	today := currentPrices[len(currentPrices)-1]
	fourDaysAgo := comparePrices[len(comparePrices)-5]

	if today < fourDaysAgo {
		return prevBuy + 1, 0
	} else if today > fourDaysAgo {
		return 0, prevSell + 1
	}
	return 0, 0
}
