package main

import (
	"fmt"

	"github.com/Antosik/rito-news/lol"
)

func Example_LoLNews(locale string, count int) {
	client := lol.NewsClient{Locale: locale}

	entries, _ := client.GetItems(count)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_LoLEsportsNews(locale string, count int) {
	client := lol.EsportsClient{Locale: locale}

	entries, _ := client.GetItems(count)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_LoLServerStatus(region string, locale string) {
	client := lol.StatusClient{Region: region}

	entries, _ := client.GetItems(locale)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func main() {
	fmt.Println("Latest League of Legends News")
	Example_LoLNews("en-us", 1)

	fmt.Println()
	fmt.Println("---")
	fmt.Println()

	fmt.Println("Latest League of Legends Esports News")
	Example_LoLEsportsNews("en-us", 1)

	fmt.Println()
	fmt.Println("---")
	fmt.Println()

	fmt.Println("Current League of Legends NA Server Status")
	Example_LoLServerStatus("na1", "en-US")
}
