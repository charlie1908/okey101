package main

import (
	"fmt"
	"okey101/Core"
	"okey101/Model"
	"strings"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
const PlayerCount int = 4
const StartTileContPerPlayer int = 21

func main() {
	playerList := make([]Model.Player, PlayerCount)

	tiles := Core.CreateFullTileSet()
	fmt.Println("Toplam taş:", len(tiles)) // 106 olmalı
	DealUserTiles(playerList, tiles)
	//var player1 = Core.ShowPlayerTiles(&tiles, "Player 1", 22)
	//var player2 = Core.ShowPlayerTiles(&tiles, "Player 2", 21)
	indicatorTile := tiles.GetRandomIndicatorFromTiles()
	fmt.Println("Indicator Tile:")
	fmt.Println(strings.Repeat("-", 30))
	colorName := Core.GetEnumName(Core.ColorEnum, indicatorTile.Color)
	fmt.Printf("ID: %d, %s %d, Joker: %v\n", indicatorTile.ID, colorName, indicatorTile.Number, indicatorTile.IsJoker)
	fmt.Println(strings.Repeat("-", 30))

	fmt.Println("Bag: ")
	for i, tile := range tiles {
		colorName := Core.GetEnumName(Core.ColorEnum, tile.Color)
		fmt.Printf("%d-) ID: %d, %s %d, Joker: %v\n", i+1, tile.ID, colorName, tile.Number, tile.IsJoker)
	}
	fmt.Println(strings.Repeat("-", 30))

	fmt.Println()
	fmt.Printf("Kalan taş sayısı: %d\n", len(tiles))

	tiles.MarkOkeyTiles(indicatorTile)
	for _, player := range playerList {
		tb := Core.TileBag(player.TileBag) // []Model.Tile → TileBag
		(&tb).MarkOkeyTiles(indicatorTile)

		fmt.Println(player.Name)
		for i, tile := range player.TileBag {
			colorName := Core.GetEnumName(Core.ColorEnum, tile.Color)
			fmt.Printf("%d-) ID: %d, %s %d, Joker: %v Okey: %v\n", i+1, tile.ID, colorName, tile.Number, tile.IsJoker, tile.IsOkey)
		}
	}
	//player1.MarkOkeyTiles(indicatorTile)
	//player2.MarkOkeyTiles(indicatorTile)

	//fmt.Println("Player 1:")
	//for i, tile := range *player1 {
	//	colorName := Core.GetEnumName(Core.ColorEnum, tile.Color)
	//	fmt.Printf("%d-) ID: %d, %s %d, Joker: %v Okey: %v\n", i+1, tile.ID, colorName, tile.Number, tile.IsJoker, tile.IsOkey)
	//}
	//for i, tile := range *player2 {
	//	colorName := Core.GetEnumName(Core.ColorEnum, tile.Color)
	//	fmt.Printf("%d-) ID: %d, %s %d, Joker: %v Okey: %v\n", i+1, tile.ID, colorName, tile.Number, tile.IsJoker, tile.IsOkey)
	//}

	fmt.Println()
	var takenTile = tiles.TakeOneFromBag(&playerList[0].TileBag)
	colorName3 := Core.GetEnumName(Core.ColorEnum, takenTile.Color)
	fmt.Printf("Ortadan Cekilen Tas - ID: %d, %s %d, Joker: %v\n", takenTile.ID, colorName3, takenTile.Number, takenTile.IsJoker)

	var dropTile = playerList[0].TileBag[3] //Player1 Bir tas Cantadan cekti ve kendine ekledi
	colorName2 := Core.GetEnumName(Core.ColorEnum, dropTile.Color)
	fmt.Printf("Player1'den Cekilen Tas - ID: %d, %s %d, Joker: %v\n", dropTile.ID, colorName2, dropTile.Number, dropTile.IsJoker)

	Core.DropTileFromTiles((*[]Model.Tile)(&playerList[0].TileBag), dropTile) //Player 1 tas cantadan ceker ve ustten 3.'yu, atar :)

	Core.TakeOneFromTable((*[]Model.Tile)(&playerList[1].TileBag), dropTile)                    //Player 2 => Player 1'in 3. elemanini ceker
	Core.DropTileFromTiles((*[]Model.Tile)(&playerList[1].TileBag), (playerList[0].TileBag)[4]) // player 2 ustten 4. elemani atar.
	for i, player := range playerList {
		start := 0
		if i == 0 {
			start = 1
		}
		tb := Core.TileBag(player.TileBag) // []Model.Tile → TileBag
		Core.ShowPlayerTiles(&tb, player.Name, StartTileContPerPlayer+start)
	}
	//Core.ShowPlayerTiles(player1, "Player 1", 22)
	//Core.ShowPlayerTiles(player2, "Player 2", 21)
	fmt.Println()
	Core.ShowPlayerTiles(&tiles, "Bag", len(tiles))

}

func DealUserTiles(playerList []Model.Player, tiles Core.TileBag) {
	for i := range PlayerCount {
		start := 0
		if i == 0 {
			start = 1
		}
		uid, err := Core.GenerateID(12)
		if err != nil {

		}
		playerList[i] = Model.Player{
			ID:       i,
			Name:     fmt.Sprintf("Player %d", i+1),
			TileBag:  *Core.ShowPlayerTiles(&tiles, "", StartTileContPerPlayer+start),
			UniqueId: uid,
		}
	}
}
