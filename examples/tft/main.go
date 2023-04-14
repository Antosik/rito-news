package main

import (
	"fmt"

	"github.com/Antosik/rito-news/tft"
)

func Example_TFTNews(locale string, count int) {
	client := tft.NewsClient{Locale: locale}

	entries, _ := client.GetItems(count)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func main() {
	fmt.Println("Latest Teamfight Tactics News")
	Example_TFTNews("en-us", 1)
}
