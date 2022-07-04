package main

import (
	"fmt"
	"os"
	lol_source "rito-news/sources/lol"
	lor_source "rito-news/sources/lor"
	riot_source "rito-news/sources/riotgames"
	tft_source "rito-news/sources/tft"
	val_source "rito-news/sources/val"
	wr_source "rito-news/sources/wr"
	"rito-news/utils"
)

func main() {
	fmt.Println("LOL")

	loln_source := lol_source.LeagueOfLegendsNews{Locale: "ru-ru"}
	loln_entries, _ := loln_source.GetItems(1)
	for _, entry := range loln_entries {
		jsonEntry, _ := utils.JSONMarshal(entry)
		os.WriteFile("loln_entries.json", jsonEntry, 0644)
	}

	fmt.Println()

	lole_source := lol_source.LeagueOfLegendsEsports{Locale: "ru-ru"}
	lole_entries, _ := lole_source.GetItems(1)
	for _, entry := range lole_entries {
		jsonEntry, _ := utils.JSONMarshal(entry)
		os.WriteFile("lole_entries.json", jsonEntry, 0644)
	}

	fmt.Println()

	lols_source := lol_source.LeagueOfLegendsStatus{Region: "br1"}
	lols_entries, _ := lols_source.GetItems("ru-RU")
	for _, entry := range lols_entries {
		jsonEntry, _ := utils.JSONMarshal(entry)
		os.WriteFile("lols_entries.json", jsonEntry, 0644)
	}

	fmt.Println()
	fmt.Println()

	fmt.Println("LOR")

	lorn_source := lor_source.LegendsOfRuneterraNews{Locale: "ru-ru"}
	lorn_entries, _ := lorn_source.GetItems(1)
	for _, entry := range lorn_entries {
		jsonEntry, _ := utils.JSONMarshal(entry)
		os.WriteFile("lorn_entries.json", jsonEntry, 0644)
	}

	fmt.Println()

	lors_source := lor_source.LegendsOfRuneterraStatus{Region: "europe"}
	lors_entries, _ := lors_source.GetItems("ru-RU")
	for _, entry := range lors_entries {
		jsonEntry, _ := utils.JSONMarshal(entry)
		os.WriteFile("lors_entries.json", jsonEntry, 0644)
	}

	fmt.Println()
	fmt.Println()

	fmt.Println("VAL")

	valn_source := val_source.VALORANTNews{Locale: "ru-ru"}
	valn_entries, _ := valn_source.GetItems(1)
	for _, entry := range valn_entries {
		jsonEntry, _ := utils.JSONMarshal(entry)
		os.WriteFile("valn_entries.json", jsonEntry, 0644)
	}

	fmt.Println()

	vale_source := val_source.VALORANTEsports{Locale: "ru-ru"}
	vale_entries, _ := vale_source.GetItems(1)
	for _, entry := range vale_entries {
		jsonEntry, _ := utils.JSONMarshal(entry)
		os.WriteFile("vale_entries.json", jsonEntry, 0644)
	}

	fmt.Println()

	vals_source := val_source.VALORANTStatus{Region: "br"}
	vals_entries, _ := vals_source.GetItems("ru-RU")
	for _, entry := range vals_entries {
		jsonEntry, _ := utils.JSONMarshal(entry)
		os.WriteFile("vals_entries.json", jsonEntry, 0644)
	}

	fmt.Println()
	fmt.Println()

	fmt.Println("WR")

	wrn_source := wr_source.WildRiftNews{Locale: "ru-ru"}
	wrn_entries, _ := wrn_source.GetItems(1)
	for _, entry := range wrn_entries {
		jsonEntry, _ := utils.JSONMarshal(entry)
		os.WriteFile("wrn_entries.json", jsonEntry, 0644)
	}

	fmt.Println()

	wre_source := wr_source.WildRiftEsports{Locale: "ru-ru"}
	wre_entries, _ := wre_source.GetItems(1)
	for _, entry := range wre_entries {
		jsonEntry, _ := utils.JSONMarshal(entry)
		os.WriteFile("wre_entries.json", jsonEntry, 0644)
	}

	fmt.Println()

	wrs_source := wr_source.WildRiftStatus{Region: "br"}
	wrs_entries, _ := wrs_source.GetItems("ru-RU")
	for _, entry := range wrs_entries {
		jsonEntry, _ := utils.JSONMarshal(entry)
		os.WriteFile("wrs_entries.json", jsonEntry, 0644)
	}

	fmt.Println()
	fmt.Println()

	fmt.Println("TFT")

	tftn_source := tft_source.TeamfightTacticsNews{Locale: "ru-ru"}
	tftn_entries, _ := tftn_source.GetItems(1)
	for _, entry := range tftn_entries {
		jsonEntry, _ := utils.JSONMarshal(entry)
		os.WriteFile("tftn_entries.json", jsonEntry, 0644)
	}

	fmt.Println()
	fmt.Println()
	fmt.Println("RITO")

	riton_source := riot_source.RiotGamesNews{Locale: "en-us"}
	riton_entries, _ := riton_source.GetItems(1)
	for _, entry := range riton_entries {
		jsonEntry, _ := utils.JSONMarshal(entry)
		os.WriteFile("riton_entries.json", jsonEntry, 0644)
	}

	fmt.Println()

	ritoj_source := riot_source.RiotGamesJobs{Locale: "en-us"}
	ritoj_entries, _ := ritoj_source.GetItems()
	for _, entry := range ritoj_entries[:1] {
		jsonEntry, _ := utils.JSONMarshal(entry)
		os.WriteFile("ritoj_entries.json", jsonEntry, 0644)
	}
}
