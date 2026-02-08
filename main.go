package main

import (
	"fmt"
	"mouniu/internal/services"
)

func main() {

	data, err := services.GetCandleStickData("hk", "01810")

	if err != nil {
		fmt.Print(err)
	}

	fmt.Print(data)

}
