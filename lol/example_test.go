package lol

import "fmt"

func Example_news() {
	client := LeagueOfLegendsNews{Locale: "ru-ru"}
	entries, _ := client.GetItems(1)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_esportsNews() {
	client := LeagueOfLegendsEsports{Locale: "ru-ru"}
	entries, _ := client.GetItems(1)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_serverStatus() {
	client := LeagueOfLegendsStatus{Region: "br1"}
	entries, _ := client.GetItems("ru-RU")
	for _, entry := range entries {
		fmt.Println(entry)
	}
}
