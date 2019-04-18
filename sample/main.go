package main

import (
	"fmt"

	"github.com/mnafees/zeit-now-go-client/now"
)

func main() {
	config := now.NewEmptyTokenConfig()
	client, _ := now.NewClient(*config)
	authToken, err := client.FetchAuthToken("<EMAIL HERE>", "Now Go Client Sample")
	if err != nil {
		panic(err)
	}
	fmt.Println(authToken)
}
