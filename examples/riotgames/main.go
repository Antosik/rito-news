package main

import (
	"fmt"

	"github.com/Antosik/rito-news/riotgames"
)

func Example_RiotGamesNews(locale string, count int) {
	client := riotgames.NewsClient{Locale: locale}

	entries, err := client.GetItems(count)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_RiotGamesJobs(locale string, count int) {
	client := riotgames.JobsClient{Locale: locale}

	entries, err := client.GetItems()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, entry := range entries[:count] {
		fmt.Println(entry)
	}
}

func main() {
	fmt.Println("Latest RiotGames News")
	Example_RiotGamesNews("en", 1)

	fmt.Println()
	fmt.Println("---")
	fmt.Println()

	fmt.Println("Some of open jobs positions at Riot Games")
	Example_RiotGamesJobs("en", 5)
}
