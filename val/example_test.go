package val

import "fmt"

func Example_news() {
	client := VALORANTNews{Locale: "ru-ru"}
	entries, _ := client.GetItems(1)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_esportsNews() {
	client := VALORANTEsports{Locale: "ru-ru"}
	entries, _ := client.GetItems(1)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func Example_serverStatus() {
	client := VALORANTStatus{Region: "br"}
	entries, _ := client.GetItems("ru-RU")
	for _, entry := range entries {
		fmt.Println(entry)
	}
}
