package tft

import "fmt"

func Example_news() {
	client := NewsClient{Locale: "ru-ru"}

	entries, _ := client.GetItems(1)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}
