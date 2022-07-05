package tft

import "fmt"

func Example_news() {
	client := TeamfightTacticsNews{Locale: "ru-ru"}
	entries, _ := client.GetItems(1)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}
