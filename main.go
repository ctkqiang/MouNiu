package main

import (
	"fmt"
	"mouniu/internal/services"
)

func main() {
	// TODO change to dynamic + grpc
	datafeed, err := services.GetCandleStickData("hk", "01810")

	if err != nil {
		fmt.Print(err)
	}

	fmt.Print(datafeed.ToJson())
}
