package main

import (
	"fmt"

	"github.com/Antosik/rito-news/val"
)

func Example_VALNews(locale string, count int) {
	client := val.NewsClient{Locale: locale}

	entries, err := client.GetItems(count)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_VALEsportsNews(locale string, count int) {
	client := val.EsportsClient{Locale: locale}

	entries, err := client.GetItems(count)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_VALServerStatus(region string, locale string) {
	client := val.StatusClient{Region: "br"}

	entries, err := client.GetItems(locale)
	if err != nil {
		fmt.Println(err)
		return
	}

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
