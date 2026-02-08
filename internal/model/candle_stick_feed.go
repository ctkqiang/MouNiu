package model

import "encoding/json"

type CandleStickData struct {
	StockName        string `json:"stock_name"`        // 股票名称
	CurrentPrice     string `json:"current_price"`     // 当前股价 [C]
	PriceChange      string `json:"price_change"`      // 涨跌额
	ChangePercentage string `json:"change_percentage"` // 涨跌幅
	PreviousClose    string `json:"previous_close"`    // 昨收盘
	TodayOpen        string `json:"today_open"`        // 今开盘 [O]
	High             string `json:"high"`              // 最高价 [H]
	Low              string `json:"low"`               // 最低价 [L]
	Volume           string `json:"volume"`            // 成交量 [V]
	Turnover         string `json:"turnover"`          // 成交额
	PERatio          string `json:"pe_ratio"`          // 市盈率
	MarketCapital    string `json:"market_capital"`    // 股市值
	Updatetime       string `json:"updatetime"`        // 更新时间
}

func (c *CandleStickData) ToJson() (string, error) {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (c *CandleStickData) ToJsonPretty() (string, error) {
	jsonData, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (c *CandleStickData) ToString() (string, error) {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
