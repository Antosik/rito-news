package main

import (
	"fmt"
	lol_source "rito-news/sources/lol"
	lor_source "rito-news/sources/lor"
	riot_source "rito-news/sources/riotgames"
	tft_source "rito-news/sources/tft"
	val_source "rito-news/sources/val"
	wr_source "rito-news/sources/wr"
)

func main() {
	fmt.Println("LOL")

	loln_source := lol_source.LeagueOfLegendsNews{Locale: "ru-ru"}
	loln_entries, _ := loln_source.GetItems(1)
	for _, entry := range loln_entries {
		fmt.Println(entry)
	}

	fmt.Println()

	lole_source := lol_source.LeagueOfLegendsEsports{Locale: "ru-ru"}
	lole_entries, _ := lole_source.GetItems(1)
	for _, entry := range lole_entries {
		fmt.Println(entry)
	}

	fmt.Println()

	lols_source := lol_source.LeagueOfLegendsStatus{Region: "br1"}
	lols_entries, _ := lols_source.GetItems("ru-RU")
	for _, entry := range lols_entries {
		fmt.Println(entry)
	}

	fmt.Println()
	fmt.Println()

	fmt.Println("LOR")

	lorn_source := lor_source.LegendsOfRuneterraNews{Locale: "ru-ru"}
	lorn_entries, _ := lorn_source.GetItems(1)
	for _, entry := range lorn_entries {
		fmt.Println(entry)
	}

	fmt.Println()

	lors_source := lor_source.LegendsOfRuneterraStatus{Region: "europe"}
	lors_entries, _ := lors_source.GetItems("ru-RU")
	for _, entry := range lors_entries {
		fmt.Println(entry)
	}

	fmt.Println()
	fmt.Println()

	fmt.Println("VAL")

	valn_source := val_source.VALORANTNews{Locale: "ru-ru"}
	valn_entries, _ := valn_source.GetItems(1)
	for _, entry := range valn_entries {
		fmt.Println(entry)
	}

	fmt.Println()

	vale_source := val_source.VALORANTEsports{Locale: "ru-ru"}
	vale_entries, _ := vale_source.GetItems(1)
	for _, entry := range vale_entries {
		fmt.Println(entry)
	}

	fmt.Println()

	vals_source := val_source.VALORANTStatus{Region: "br"}
	vals_entries, _ := vals_source.GetItems("ru-RU")
	for _, entry := range vals_entries {
		fmt.Println(entry)
	}

	fmt.Println()
	fmt.Println()

	fmt.Println("WR")

	wrn_source := wr_source.WildRiftNews{Locale: "ru-ru"}
	wrn_entries, _ := wrn_source.GetItems(1)
	for _, entry := range wrn_entries {
		fmt.Println(entry)
	}

	fmt.Println()

	wre_source := wr_source.WildRiftEsports{Locale: "ru-ru"}
	wre_entries, _ := wre_source.GetItems(1)
	for _, entry := range wre_entries {
		fmt.Println(entry)
	}

	fmt.Println()

	wrs_source := wr_source.WildRiftStatus{Region: "br"}
	wrs_entries, _ := wrs_source.GetItems("ru-RU")
	for _, entry := range wrs_entries {
		fmt.Println(entry)
	}

	fmt.Println()
	fmt.Println()

	fmt.Println("TFT")

	tftn_source := tft_source.TeamfightTacticsNews{Locale: "ru-ru"}
	tftn_entries, _ := tftn_source.GetItems(1)
	for _, entry := range tftn_entries {
		fmt.Println(entry)
	}

	fmt.Println()
	fmt.Println()
	fmt.Println("RITO")

	riton_source := riot_source.RiotGamesNews{Locale: "en-us"}
	riton_entries, _ := riton_source.GetItems(1)
	for _, entry := range riton_entries {
		fmt.Println(entry)
	}

	fmt.Println()

	ritoj_source := riot_source.RiotGamesJobs{Locale: "en-us"}
	ritoj_entries, _ := ritoj_source.GetItems()
	for _, entry := range ritoj_entries[:1] {
		fmt.Println(entry)
	}
}
