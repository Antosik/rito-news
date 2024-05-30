package main

import (
	"fmt"

	"github.com/Antosik/rito-news/tft"
)

func Example_TFTNews(locale string, count int) {
	client := tft.NewsClient{Locale: locale}

	entries, err := client.GetItems(count)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func main() {
	fmt.Println("Latest Teamfight Tactics News")
	Example_TFTNews("en-us", 1)
}
