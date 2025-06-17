package Core

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"okey101/Model"
	"strings"
)

// OkeyTaşınıBelirle, gösterge taşına göre gerçek okey taşını döner
func DetermineOkeyTile(indicator Model.Tile) Model.Tile {
	nextNumber := indicator.Number + 1
	if nextNumber > 13 {
		nextNumber = 1
	}

	return Model.Tile{
		Number:  nextNumber,
		Color:   indicator.Color,
		IsOkey:  true,
		IsJoker: false,
	}
}

/*// MarkOkeyTiles, göstergeye göre okey olan taşları işaretler
func (tiles *TileBag) MarkOkeyJokerTiles(indicator Model.Tile) {
	okey := DetermineOkeyTile(indicator)
	for i := range *tiles {
		tile := &(*tiles)[i]
		//Set Fake Joker
		if tile.Number == indicator.Number && tile.Color == indicator.Color {
			tile.IsJoker = true
			//Set Okey
		} else if tile.Number == okey.Number && tile.Color == okey.Color && !tile.IsJoker {
			tile.IsOkey = true
		}
	}
}*/

// MarkOkeyTiles, göstergeye göre okey olan taşları işaretler
/*func (tiles *TileBag) MarkOkeyTiles(indicator Model.Tile) {
	okey := DetermineOkeyTile(indicator)
	for i := range *tiles {
		tile := &(*tiles)[i]
		//Set Okey
		if tile.Number == okey.Number && tile.Color == okey.Color && !tile.IsJoker {
			tile.IsOkey = true
			//Set Okey
		} else {
			tile.IsOkey = false
		}
	}
}*/

/*func (tiles *TileBag) MarkOkeyTiles(indicator Model.Tile) {
	okey := DetermineOkeyTile(indicator)
	for i := range *tiles {
		tile := &(*tiles)[i]
		if tile.Number == indicator.Number && tile.Color == indicator.Color {
			tile.IsJoker = true // Sahte okey - Gosterge
			tile.IsOkey = false
		} else if tile.Number == okey.Number && tile.Color == okey.Color && !tile.IsJoker {
			tile.IsOkey = true // Gerçek okey
		} else {
			tile.IsOkey = false
			tile.IsJoker = false
		}
	}
}*/

func (tiles *TileBag) MarkOkeyTiles(indicator Model.Tile) {
	okey := DetermineOkeyTile(indicator)

	for i := range *tiles {
		tile := &(*tiles)[i]

		// GOSTERGE TASI ICIN NE YAPMAK GEREKIR ?
		if tile.Number == indicator.Number && tile.Color == indicator.Color {
			// Gösterge taşı -> sahte okey olarak işaretle ama numara ve renk belirtme
			tile.IsJoker = true
			tile.IsOkey = false
		} else if tile.Number == okey.Number && tile.Color == okey.Color && !tile.IsJoker {
			// Gerçek okey
			tile.IsOkey = true
		} else if tile.Number == 0 && tile.Color == ColorEnum.None && tile.IsJoker {
			// Sahte okey taşları — göstergeye göre joker gibi davranmalı
			tile.Number = okey.Number
			tile.Color = okey.Color
		} else {
			// Normal taş
			tile.IsOkey = false
			tile.IsJoker = false
		}
	}
}

func CreateFullTileSet() TileBag {
	var tiles []Model.Tile
	id := 0

	colors := []int{
		ColorEnum.Red,
		ColorEnum.Yellow,
		ColorEnum.Blue,
		ColorEnum.Black,
	}

	for _, color := range colors {
		for number := 1; number <= 13; number++ {
			for i := 0; i < 2; i++ { // Her taştan 2 adet
				tiles = append(tiles, Model.Tile{
					ID:      id,
					Number:  number,
					Color:   color,
					IsJoker: false,
					IsOkey:  false,
					IsOpend: false,
				})
				id++
			}
		}
	}

	// 2 adet sahte okey taşı (renksiz, numbersız)
	for i := 0; i < 2; i++ {
		tiles = append(tiles, Model.Tile{
			ID:      id,
			Number:  0,
			Color:   ColorEnum.None,
			IsJoker: true,
			IsOkey:  false,
		})
		id++
	}

	return ShuffleTilesSecure(tiles)
}

// Fisher-Yates Shuffle :)
// ShuffleTilesSecure, taşları kriptografik olarak güvenli şekilde karıştırır
func ShuffleTilesSecure(tiles []Model.Tile) TileBag {
	shuffled := make([]Model.Tile, len(tiles))
	copy(shuffled, tiles)

	for i := len(shuffled) - 1; i > 0; i-- {
		j, err := cryptoRandInt(i + 1)
		if err != nil {
			// Hata durumunda karıştırmadan döner
			return tiles
		}
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled
}

// cryptoRandInt returns a random int between 0 and max-1 using crypto/rand
func cryptoRandInt(max int) (int, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(nBig.Int64()), nil
}

type TileBag []Model.Tile

func (tiles *TileBag) GetTiles(count int) *TileBag {
	// örnek implementasyon
	if len(*tiles) < count {
		count = len(*tiles)
	}

	selected := (*tiles)[:count]
	*tiles = (*tiles)[count:] // kalan taşlar

	return &selected
}

func ShowPlayerTiles(tiles *TileBag, name string, topCount int) *TileBag {
	player := tiles.GetTiles(topCount)

	fmt.Println(name)
	fmt.Println(strings.Repeat("-", 30))

	for i, tile := range *player {
		colorName := GetEnumName(ColorEnum, tile.Color)
		//fmt.Printf("%d-) ID: %d, %s %d, Joker: %v\n", i+1, tile.ID, colorName, tile.Number, tile.IsJoker)
		fmt.Printf("%d-) ID: %d, %s %d, Joker: %v Okey: %v\n", i+1, tile.ID, colorName, tile.Number, tile.IsJoker, tile.IsOkey)
	}

	fmt.Println()
	fmt.Printf("Kalan taş sayısı: %d\n", len(*tiles))
	return player
}

// Bunun yerine en usttekini de alabiliriz. Ben Random secmeyi daha guvenli buldum.
func (tiles *TileBag) GetRandomIndicatorFromTiles() Model.Tile {
	// Joker olmayan taşları filtrele
	validTiles := make([]Model.Tile, 0)
	for _, tile := range *tiles {
		if !tile.IsJoker {
			validTiles = append(validTiles, tile)
		}
	}

	if len(validTiles) == 0 {
		log.Fatal("TileBag contains no valid (non-Joker) tiles for indicator selection")
	}

	// Güvenli rastgele index seç
	max := big.NewInt(int64(len(validTiles)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatalf("crypto/rand failed: %v", err)
	}
	randomIndex := int(n.Int64())
	indicator := validTiles[randomIndex]

	// Seçilen taşı orijinal TileBag'den çıkar
	DropTileFromTiles((*[]Model.Tile)(tiles), indicator)
	/*for i, tile := range *tiles {
		if tile.ID == indicator.ID {
			*tiles = append((*tiles)[:i], (*tiles)[i+1:]...)
			break
		}
	}*/

	return indicator
}

func (tiles *TileBag) TakeOneFromBag(player *[]Model.Tile) Model.Tile {
	if len(*tiles) == 0 {
		log.Fatal("TileBag contains no valid (non-Joker) tiles for indicator selection")
	}
	var tile = (*tiles)[0]
	*player = append(*player, tile)
	return tile
}

func TakeOneFromTable(player *[]Model.Tile, tile Model.Tile) {
	*player = append(*player, tile)
}

/*func DropTileFromTiles(playerTiles *[]Model.Tile, dropTile Model.Tile) {
	for i, tile := range *playerTiles {
		if tile.ID == dropTile.ID {
			*playerTiles = append((*playerTiles)[:i], (*playerTiles)[i+1:]...)
			break
		}
	}
}*/

func DropTileFromTiles(playerTiles *[]Model.Tile, dropTile Model.Tile) bool {
	var isFound bool = false
	for i, tile := range *playerTiles {
		if tile.ID == dropTile.ID {
			*playerTiles = append((*playerTiles)[:i], (*playerTiles)[i+1:]...)
			isFound = true
			break
		}
	}
	return isFound
}

func FloatPtr(f float64) *float64 {
	return &f
}

func IntPtr(i int) *int {
	return &i
}

func BoolPtr(b bool) *bool {
	return &b
}

func StringPtr(s string) *string { return &s }

var Game GameGroupState

func ResetGame() {
	Game = GameGroupState{
		GroupIDCounter: 0,
	}
}

type GameGroupState struct {
	GroupIDCounter int
}

func (g *GameGroupState) GenerateGroupID() int {
	g.GroupIDCounter++
	return g.GroupIDCounter
}
