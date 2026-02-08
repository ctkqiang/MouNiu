# -*- coding: UTF-8 -*-

try:
    import sys
    import akshare as ak
except ImportError as e:
    raise ImportError(f"{e}  请使用 pip install -r requirements.txt")

class StockData:
    def __init__(self, tickerId: str) -> None:
        self.tickerId = tickerId
        print(f"AkShare 版本: {ak.__version__}")
        print(f"正在处理股票代码: {self.tickerId}")

    def get_stock_hist(self, symbolID :str) -> None:
        pass

if __name__ == "__main__":
    ticker = None
    args = sys.argv[1:]

    for arg in args:
        if arg.startswith("ticker="):
            ticker = arg.split("=", 1)[1]
        
    if ticker is None:
        print("用法: python3 scripts/run.py ticker=XXXXXX.XX")
        print("示例: python3 scripts/run.py ticker=01810.HK")
        sys.exit(1)

    stock_data = StockData(tickerId=ticker)
