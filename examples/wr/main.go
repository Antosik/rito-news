package main

import (
	"fmt"

	"github.com/Antosik/rito-news/wr"
)

func Example_WRNews(locale string, count int) {
	client := wr.NewsClient{Locale: locale}

	entries, _ := client.GetItems(count)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_WREsportsNews(locale string, count int) {
	client := wr.EsportsClient{Locale: locale}

	entries, _ := client.GetItems(count)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_WRServerStatus(region string, locale string) {
	client := wr.StatusClient{Region: region}

	entries, _ := client.GetItems(locale)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func main() {
	fmt.Println("Latest League of Legends: Wild Rift News")
	Example_WRNews("en-us", 1)

	fmt.Println()
	fmt.Println("---")
	fmt.Println()

	fmt.Println("Latest League of Legends: Wild Rift Esports News")
	Example_WREsportsNews("en-us", 1)

	fmt.Println()
	fmt.Println("---")
	fmt.Println()

	fmt.Println("Current League of Legends: Wild Rift NA Server Status")
	Example_WRServerStatus("na", "en-US")
}
