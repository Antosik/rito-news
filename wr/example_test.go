package wr

import "fmt"

func Example_news() {
	client := NewsClient{Locale: "ru-ru"}

	entries, _ := client.GetItems(1)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_esportsNews() {
	client := EsportsClient{Locale: "ru-ru"}

	entries, _ := client.GetItems(1)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_serverStatus() {
	client := StatusClient{Region: "br1"}

	entries, _ := client.GetItems("ru-RU")
	for _, entry := range entries {
		fmt.Println(entry)
	}
}
