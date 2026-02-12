package config

var (
	API string = "/api"

	ANNOUNCEMENT_ALL     string = "/announcement/all"
	ANNOUNCEMENT_SYMBOLS string = "/announcement/:symbol"

	STOCKS_ALL                  string = "/stocks/all"
	STOCKS_SYMBOL               string = "/:ticker"
	STOCKS_SYMBOL_CURRENT_PRICE string = "/current/:ticker"

	ANALYSIS_ALL string = "/analysis"
)
