package test

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func GetStockDetails() (map[string]string, error) {
	result := make(map[string]string)

	// 使用新浪财经港股API接口
	url := "https://hq.sinajs.cn/list=hk01810"

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	// 设置完整的浏览器User-Agent头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://finance.sina.com.cn")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取API响应
	reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	body := string(bodyBytes)

	// 解析新浪财经API格式: var hq_str_hk01810="小米集团-W,12.340,+0.20,+1.65%,122.45,122.45,122.20,122.60,3144000,385209600,34.22,250亿,2024-02-08 16:08:00,0";
	if strings.Contains(body, "var hq_str_hk01810=") {
		// 提取数据部分
		start := strings.Index(body, "\"")
		end := strings.LastIndex(body, "\"")
		if start != -1 && end != -1 && end > start {
			dataStr := body[start+1 : end]
			fields := strings.Split(dataStr, ",")

			if len(fields) >= 14 {
				result["name"] = fields[0]
				result["price"] = fields[1]
				result["change"] = fields[2]
				result["change_pct"] = fields[3]
				result["prev_close"] = fields[4]
				result["open"] = fields[5]
				result["high"] = fields[6]
				result["low"] = fields[7]
				result["volume"] = fields[8] + "股"
				result["turnover"] = fields[9] + "元"
				result["pe_ratio"] = fields[10]
				result["market_cap"] = fields[11]
				result["update_time"] = fields[12]
			}
		}
		return result, nil
	}

	// 如果API失败，回退到HTML解析
	return parseHTMLData()
}

func parseHTMLData() (map[string]string, error) {
	result := make(map[string]string)

	// 使用HTML页面作为备选
	url := "https://stock.finance.sina.com.cn/hkstock/quotes/01810.html"

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 处理GBK编码转换为UTF-8
	utf8Reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())

	doc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败: %v", err)
	}

	// 提取所有文本内容进行正则匹配
	allText := doc.Text()

	// 使用更精确的选择器提取数据
	// 提取股票名称
	name := doc.Find("#stock_cname").Text()
	if name == "" {
		name = doc.Find(".name01").Text()
	}
	result["name"] = strings.TrimSpace(name)

	// 提取当前价格 - 尝试多个可能的选择器
	price := doc.Find("#mts_stock_hk_price").Text()
	if price == "" || price == "--" {
		price = doc.Find(".price.fl").Text()
	}
	if price == "" || price == "--" {
		// 从标题中提取
		title := doc.Find("title").Text()
		if match := regexp.MustCompile(`(\d+\.\d+)`).FindStringSubmatch(title); len(match) > 0 {
			price = match[0]
		}
	}
	result["price"] = strings.TrimSpace(price)

	// 提取涨跌额和涨跌幅
	change := doc.Find("#mts_stock_hk_zdf").Text()
	if change != "" && change != "--(--.-%)" {
		// 解析格式如: +0.20(+0.16%)
		if match := regexp.MustCompile(`([+-]?\d+\.\d+)\(([+-]?\d+\.\d+%)\)`).FindStringSubmatch(change); len(match) > 2 {
			result["change"] = strings.TrimSpace(match[1])
			result["change_pct"] = strings.TrimSpace(match[2])
		}
	}

	// 提取主要行情数据 - 更精确的OHLCV模式
	patterns := map[string]string{
		"prev_close":   `昨收盘\s*([\d.]+)`,
		"high":         `最高价\s*([\d.]+)`,
		"low":          `最低价\s*([\d.]+)`,
		"open":         `今开盘\s*([\d.]+)`,
		"change":       `涨跌额\s*([+-]?[\d.]+)`,
		"change_pct":   `涨跌幅\s*([+-]?[\d.]+%)`,
		"volume":       `成交量\s*([\d.]+[万亿]?)`,
		"turnover":     `成交额\s*([\d.]+[万亿]?)`,
		"pe_ratio":     `市盈率\s*([\d.]+)`,
		"market_cap":   `港股市值\s*([\d.]+[万亿]?)`,
		"lot_size":     `每手股数\s*(\d+)`,
		"week_52_high": `52周最高\s*([\d.]+)`,
		"week_52_low":  `52周最低\s*([\d.]+)`,
	}

	for key, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if match := re.FindStringSubmatch(allText); len(match) > 1 {
			result[key] = strings.TrimSpace(match[1])
		}
	}

	// 专门处理OHLCV数据 - 从表格中提取
	result = extractOHLCVData(doc, result, allText)

	// 提取更新时间
	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2})`)
	if match := re.FindStringSubmatch(allText); len(match) > 0 {
		result["update_time"] = match[1]
	}

	return result, nil
}

// 辅助函数：读取GBK编码的响应
func readGBKResponse(body io.ReadCloser) (string, error) {
	defer body.Close()

	// 创建GBK到UTF-8的转换器
	reader := transform.NewReader(body, simplifiedchinese.GBK.NewDecoder())
	scanner := bufio.NewScanner(reader)
	var content strings.Builder

	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return content.String(), nil
}

// extractOHLCVData 专门提取OHLCV数据
func extractOHLCVData(doc *goquery.Document, result map[string]string, allText string) map[string]string {
	// 方法1: 从表格中提取OHLCV数据
	doc.Find("table tr").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(text, "开盘") || strings.Contains(text, "最高") || strings.Contains(text, "最低") || strings.Contains(text, "收盘") {
			// 提取表格数据
			s.Find("td").Each(func(j int, td *goquery.Selection) {
				tdText := strings.TrimSpace(td.Text())
				switch j {
				case 0:
					if strings.Contains(tdText, "开盘") && result["open"] == "" {
						result["open"] = extractNumberFromText(tdText)
					}
				case 1:
					if strings.Contains(tdText, "最高") && result["high"] == "" {
						result["high"] = extractNumberFromText(tdText)
					}
				case 2:
					if strings.Contains(tdText, "最低") && result["low"] == "" {
						result["low"] = extractNumberFromText(tdText)
					}
				case 3:
					if strings.Contains(tdText, "收盘") && result["prev_close"] == "" {
						result["prev_close"] = extractNumberFromText(tdText)
					}
				}
			})
		}
	})

	// 方法2: 使用正则表达式从文本中提取OHLCV
	if result["open"] == "" {
		if match := regexp.MustCompile(`今开盘\s*([\d.]+)`).FindStringSubmatch(allText); len(match) > 1 {
			result["open"] = match[1]
		} else if match := regexp.MustCompile(`开盘[\s\S]*?(\d+\.\d+)`).FindStringSubmatch(allText); len(match) > 1 {
			result["open"] = match[1]
		}
	}

	if result["high"] == "" {
		if match := regexp.MustCompile(`最高价\s*([\d.]+)`).FindStringSubmatch(allText); len(match) > 1 {
			result["high"] = match[1]
		} else if match := regexp.MustCompile(`最高[\s\S]*?(\d+\.\d+)`).FindStringSubmatch(allText); len(match) > 1 {
			result["high"] = match[1]
		} else if match := regexp.MustCompile(`高[\s\S]*?(\d+\.\d+)`).FindStringSubmatch(allText); len(match) > 1 {
			result["high"] = match[1]
		}
	}

	if result["low"] == "" {
		if match := regexp.MustCompile(`最低价\s*([\d.]+)`).FindStringSubmatch(allText); len(match) > 1 {
			result["low"] = match[1]
		} else if match := regexp.MustCompile(`最低[\s\S]*?(\d+\.\d+)`).FindStringSubmatch(allText); len(match) > 1 {
			result["low"] = match[1]
		} else if match := regexp.MustCompile(`低[\s\S]*?(\d+\.\d+)`).FindStringSubmatch(allText); len(match) > 1 {
			result["low"] = match[1]
		}
	}

	if result["prev_close"] == "" {
		if match := regexp.MustCompile(`昨收盘\s*([\d.]+)`).FindStringSubmatch(allText); len(match) > 1 {
			result["prev_close"] = match[1]
		} else if match := regexp.MustCompile(`收盘[\s\S]*?(\d+\.\d+)`).FindStringSubmatch(allText); len(match) > 1 {
			result["prev_close"] = match[1]
		}
	}

	// 确保收盘价就是当前价格
	if result["prev_close"] == "" && result["price"] != "" {
		result["prev_close"] = result["price"]
	}

	return result
}

// extractNumberFromText 从文本中提取数字
func extractNumberFromText(text string) string {
	if match := regexp.MustCompile(`([\d.]+)`).FindStringSubmatch(text); len(match) > 1 {
		return match[1]
	}
	return ""
}

// data, err := services.GetStockDetails()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("股票名称: %s\n", data["name"])
// 	fmt.Printf("当前股价: %s\n", data["price"])
// 	fmt.Printf("涨跌额: %s\n", data["change"])
// 	fmt.Printf("涨跌幅: %s\n", data["change_pct"])
// 	fmt.Printf("昨收盘: %s\n", data["prev_close"])
// 	fmt.Printf("今开盘: %s\n", data["open"])
// 	fmt.Printf("最高价: %s\n", data["high"])
// 	fmt.Printf("最低价: %s\n", data["low"])
// 	fmt.Printf("成交量: %s\n", data["volume"])
// 	fmt.Printf("成交额: %s\n", data["turnover"])
// 	fmt.Printf("市盈率: %s\n", data["pe_ratio"])
// 	fmt.Printf("港股市值: %s\n", data["market_cap"])
// 	fmt.Printf("更新时间: %s\n", data["update_time"])
