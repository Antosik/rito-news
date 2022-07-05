package wr

import "fmt"

func Example_news() {
	client := WildRiftNews{Locale: "ru-ru"}
	entries, _ := client.GetItems(1)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_esportsNews() {
	client := WildRiftEsports{Locale: "ru-ru"}
	entries, _ := client.GetItems(1)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_serverStatus() {
	client := WildRiftStatus{Region: "br1"}
	entries, _ := client.GetItems("ru-RU")
	for _, entry := range entries {
		fmt.Println(entry)
	}
}
