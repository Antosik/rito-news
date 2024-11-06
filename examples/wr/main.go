package main

import (
	"fmt"

	"github.com/Antosik/rito-news/wr"
)

func Example_WRNews(locale string, count int) {
	client := wr.NewsClient{Locale: locale}

	entries, err := client.GetItems(count)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_WRServerStatus(region string, locale string) {
	client := wr.StatusClient{Region: region}

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
	fmt.Println("Latest League of Legends: Wild Rift News")
	Example_WRNews("en-us", 1)

	fmt.Println()
	fmt.Println("---")
	fmt.Println()

	fmt.Println("Current League of Legends: Wild Rift NA Server Status")
	Example_WRServerStatus("na", "en-US")
}
