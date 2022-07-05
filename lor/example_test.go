package lor

import "fmt"

func Example_news() {
	client := LegendsOfRuneterraNews{Locale: "ru-ru"}
	entries, _ := client.GetItems(1)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_serverStatus() {
	client := LegendsOfRuneterraStatus{Region: "europe"}
	entries, _ := client.GetItems("ru-RU")
	for _, entry := range entries {
		fmt.Println(entry)
	}
}
