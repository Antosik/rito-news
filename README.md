# rito-news

[![Go Reference](https://pkg.go.dev/badge/github.com/Antosik/rito-news.svg)](https://pkg.go.dev/github.com/Antosik/rito-news) [![Go Report Card](https://goreportcard.com/badge/github.com/Antosik/rito-news)](https://goreportcard.com/report/github.com/Antosik/rito-news)

## Description

Go package that provides an API to get official news about [Riot Games](https://www.riotgames.com) and their games ([League of Legends](https://leagueoflegends.com/), [Legends of Runeterra](https://playruneterra.com/), [Teamfight Tactics](https://teamfighttactics.leagueoflegends.com), [VALORANT](https://playvalorant.com/) and [Wild Rift](https://wildrift.leagueoflegends.com/))

## How to use

```go
package main

import (
	"fmt"

	"github.com/Antosik/rito-news/lol"
)

func main() {
	client := lol.NewsClient{Locale: "ru-ru"}
	items, err := client.GetItems(10)

	if err != nil {
		fmt.Printf("An error occured: %v", err)
	} else {
		for _, item := range items {
			fmt.Println(item.Title, item.URL)
		}
	}
}
```

## Supported Services

-   League of Legends ([examples](https://github.com/Antosik/rito-news/blob/main/examples/lol/main.go))
    -   [News](https://www.leagueoflegends.com/en-us/news/)
    -   [Esports](https://lolesports.com/news)
    -   [Server status](https://status.riotgames.com/lol?region=na1&locale=en_US)
-   Legends of Runeterra ([examples](https://github.com/Antosik/rito-news/blob/main/examples/lor/main.go))
    -   [News](https://playruneterra.com/en-us/news/)
    -   [Server status](https://status.riotgames.com/lor?region=europe&locale=en_US)
-   Riot Games ([examples](https://github.com/Antosik/rito-news/blob/main/examples/riotgames/main.go))
    -   [News](https://www.riotgames.com/en/news)
    -   [Jobs](https://www.riotgames.com/en/work-with-us)
-   Teamfight Tactics ([examples](https://github.com/Antosik/rito-news/tree/main/examples/tft/main.go))
    -   [News](https://teamfighttactics.leagueoflegends.com/en-us/news/)
-   VALORANT ([examples](https://github.com/Antosik/rito-news/tree/main/examples/val/main.go))
    -   [News](https://playvalorant.com/en-us/news/)
    -   [Esports](https://valorantesports.com/news)
    -   [Server status](https://status.riotgames.com/valorant?region=na&locale=en_US)
-   Wild Rift ([examples](https://github.com/Antosik/rito-news/blob/main/examples/wr/main.go))
    -   [News](https://wildrift.leagueoflegends.com/en-us/news/)
    -   [Esports](https://wildriftesports.com/en-us/news)
    -   [Server status](https://status.riotgames.com/wildrift?region=na&locale=en_US)

## Attribution
This service isn't developed by Riot Games and doesn't reflect the views or opinions of Riot Games or anyone officially involved in producing or managing League of Legends, Legends of Runeterra, Teamfight Tactics, VALORANT, or Wild Rift. League of Legends, Legends of Runeterra, Teamfight Tactics, VALORANT, Wild Rift and Riot Games are trademarks or registered trademarks of Riot Games, Inc. League of Legends, Legends of Runeterra, Teamfight Tactics, VALORANT, Wild Rift (c) Riot Games, Inc.
