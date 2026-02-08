package services

import (
	"mouniu/internal/utilities"
	"os/exec"
)

func GetCandlestickData(tickerId string) {
	cmd := exec.Command("python3", "../../scripts/run.py", "ticker="+tickerId)
	output, err := cmd.CombinedOutput()

	if err != nil {
		utilities.Error("执行脚本失败: %v, 输出: %s", err, string(output))
		return
	}

	utilities.Info("脚本执行成功，输出如下:\n%s", string(output))
}
