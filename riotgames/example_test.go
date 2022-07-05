package riotgames

import "fmt"

func Example_news() {
	client := NewsClient{Locale: "ru-ru"}

	entries, _ := client.GetItems(1)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_jobs() {
	client := JobsClient{Locale: "ru-ru"}

	entries, _ := client.GetItems()
	for _, entry := range entries {
		fmt.Println(entry)
	}
}
