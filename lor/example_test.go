package lor

import "fmt"

func Example_news() {
	client := NewsClient{Locale: "ru-ru"}

	entries, _ := client.GetItems(1)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_serverStatus() {
	client := StatusClient{Region: "europe"}

	entries, _ := client.GetItems("ru-RU")
	for _, entry := range entries {
		fmt.Println(entry)
	}
}
