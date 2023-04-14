package main

import (
	"fmt"

	"github.com/Antosik/rito-news/lor"
)

func Example_LoRNews(locale string, count int) {
	client := lor.NewsClient{Locale: locale}

	entries, _ := client.GetItems(count)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_LoRServerStatus(region string, locale string) {
	client := lor.StatusClient{Region: region}

	entries, _ := client.GetItems(locale)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func main() {
	fmt.Println("Latest Legends of Runeterra News")
	Example_LoRNews("en-us", 1)

	fmt.Println()
	fmt.Println("---")
	fmt.Println()

	fmt.Println("Current Legends of Runeterra Americas Server Status")
	Example_LoRServerStatus("americas", "en-US")
}
