package model

// MACDResult 包含“指数平滑异同移动平均线” (MACD) 指标的计算结果
type MACDResult struct {
	DIF  float64 `json:"dif"`  // 离差值 (EMA12 - EMA26)
	DEA  float64 `json:"dea"`  // 讯号线 (DIF 的 9 日 EMA)
	MACD float64 `json:"macd"` // 柱状图 (DIF - DEA) * 2
}

// BollingerResult 包含“布林带” (Bollinger Bands) 指标的计算结果
type BollingerResult struct {
	Upper  float64 `json:"upper"`  // 上轨线 (压力线)
	Middle float64 `json:"middle"` // 中轨线 (通常是 20 日 SMA)
	Lower  float64 `json:"lower"`  // 下轨线 (支撑线)
}

// TD9Result 包含“神奇九转” (TD Sequential) 指标的计算结果
type TD9Result struct {
	Count       int  `json:"count"`         // 当前计数 (1-9)，表示连续满足条件的周期数
	IsBuySetup  bool `json:"is_buy_setup"`  // 是否触发买入结构 (预示超跌反弹)
	IsSellSetup bool `json:"is_sell_setup"` // 是否触发卖出结构 (预示超买回调)
}
