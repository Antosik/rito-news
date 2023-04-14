package main

import (
	"fmt"

	"github.com/Antosik/rito-news/val"
)

func Example_VALNews(locale string, count int) {
	client := val.NewsClient{Locale: locale}

	entries, _ := client.GetItems(count)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_VALEsportsNews(locale string, count int) {
	client := val.EsportsClient{Locale: locale}

	entries, _ := client.GetItems(count)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_VALServerStatus(region string, locale string) {
	client := val.StatusClient{Region: "br"}

	entries, _ := client.GetItems(locale)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func main() {
	fmt.Println("Latest VALORANT News")
	Example_VALNews("en-us", 1)

	fmt.Println()
	fmt.Println("---")
	fmt.Println()

	fmt.Println("Latest VALORANT Esports News")
	Example_VALEsportsNews("en-us", 1)

	fmt.Println()
	fmt.Println("---")
	fmt.Println()

	fmt.Println("Current VALORANT NA Server Status")
	Example_VALServerStatus("na", "en-US")
}
