package main

import (
	"fmt"
	"okey101/Core"
	"okey101/Model"
	"strings"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	tiles := Core.CreateFullTileSet()
	fmt.Println("Toplam taş:", len(tiles)) // 106 olmalı
	var player1 = Core.ShowPlayerTiles(&tiles, "Player 1:", 22)
	var player2 = Core.ShowPlayerTiles(&tiles, "Player 2:", 21)
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
	player1.MarkOkeyTiles(indicatorTile)
	player2.MarkOkeyTiles(indicatorTile)

	fmt.Println("Player 1:")
	for i, tile := range *player1 {
		colorName := Core.GetEnumName(Core.ColorEnum, tile.Color)
		fmt.Printf("%d-) ID: %d, %s %d, Joker: %v Okey: %v\n", i+1, tile.ID, colorName, tile.Number, tile.IsJoker, tile.IsOkey)
	}
	for i, tile := range *player2 {
		colorName := Core.GetEnumName(Core.ColorEnum, tile.Color)
		fmt.Printf("%d-) ID: %d, %s %d, Joker: %v Okey: %v\n", i+1, tile.ID, colorName, tile.Number, tile.IsJoker, tile.IsOkey)
	}

	fmt.Println()
	var takenTile = tiles.TakeOneFromBag((*[]Model.Tile)(player1))
	colorName3 := Core.GetEnumName(Core.ColorEnum, takenTile.Color)
	fmt.Printf("Ortadan Cekilen Tas - ID: %d, %s %d, Joker: %v\n", takenTile.ID, colorName3, takenTile.Number, takenTile.IsJoker)

	var dropTile = (*player1)[3] //Player1 Bir tas Cantadan cekti ve kendine ekledi
	colorName2 := Core.GetEnumName(Core.ColorEnum, dropTile.Color)
	fmt.Printf("Player1'den Cekilen Tas - ID: %d, %s %d, Joker: %v\n", dropTile.ID, colorName2, dropTile.Number, dropTile.IsJoker)

	Core.DropTileFromTiles((*[]Model.Tile)(player1), dropTile) //Player 1 tas cantadan ceker ve ustten 3.'yu, atar :)

	Core.TakeOneFromTable((*[]Model.Tile)(player2), dropTile)       //Player 2 => Player 1'in 3. elemanini ceker
	Core.DropTileFromTiles((*[]Model.Tile)(player2), (*player2)[4]) // player 2 ustten 4. elemani atar.

	Core.ShowPlayerTiles(player1, "Player 1:", 22)
	Core.ShowPlayerTiles(player2, "Player 2:", 21)
	fmt.Println()
	Core.ShowPlayerTiles(&tiles, "Bag :", len(tiles))

}
