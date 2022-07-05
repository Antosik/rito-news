package riotgames

import "fmt"

func Example_news() {
	client := RiotGamesNews{Locale: "ru-ru"}
	entries, _ := client.GetItems(1)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_jobs() {
	client := RiotGamesJobs{Locale: "ru-ru"}
	entries, _ := client.GetItems()
	for _, entry := range entries {
		fmt.Println(entry)
	}
}
