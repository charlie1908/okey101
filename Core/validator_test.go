package Core

import (
	"okey101/Model"
	"testing"
)

// Test Group gecerli olmali.
// Gecersiz yani false oldugu durumda Test yanlistir.
func TestValidGroup(t *testing.T) {
	tiles := []Model.Tile{
		{Number: 5, Color: 1},
		{Number: 5, Color: 2},
		{Number: 5, Color: 3},
	}
	if !IsValidGroupOrRun(tiles) {
		t.Error("Expected valid group, got invalid")
	} else {
		t.Log("PASS Valid Group")
	}
}

// Bu mutlaka false donmeli. True donerse valid gozukur. Test hatali demektir.
// **Test Group gecersiz oldugu zaman dogru ve gecerlidir.
func TestInvalidGroup_DuplicateColor(t *testing.T) {
	tiles := []Model.Tile{
		{Number: 5, Color: 1},
		{Number: 5, Color: 1},
		{Number: 5, Color: 2},
	}
	if IsValidGroupOrRun(tiles) {
		t.Error("Expected invalid group due to duplicate color")
	} else {
		t.Log("PASS Invalid Group (duplicate color)")
	}
}

// Sirali Artan ayni Renk Run Gecerli Testi.
func TestValidRun(t *testing.T) {
	tiles := []Model.Tile{
		{Number: 4, Color: 1},
		{Number: 5, Color: 1},
		{Number: 6, Color: 1},
	}
	if !IsValidGroupOrRun(tiles) {
		t.Error("Expected valid run, got invalid")
	} else {
		t.Log("PASS Valid Run")
	}
}

// Sirali Artan ayni Renk Okeyli Run Gecerli Testi.
func TestValidRunWithOkey(t *testing.T) {
	tiles := []Model.Tile{
		{Number: 4, Color: 1},
		{IsOkey: true},
		{IsOkey: true},
		{Number: 7, Color: 1},
	}
	if !IsValidGroupOrRun(tiles) {
		t.Error("Expected valid run with joker, got invalid")
	} else {
		t.Log("PASS Valid Run with Okey")
	}
}

func TestCanOpenTiles_Valid(t *testing.T) {
	opened := [][]Model.Tile{
		{
			{Number: 5, Color: 1},
			{Number: 5, Color: 2},
			{Number: 5, Color: 3},
		},
		{
			{Number: 11, Color: 1},
			{Number: 12, Color: 1},
			{Number: 13, Color: 1},
		},
		{
			{Number: 6, Color: 2},
			{Number: 6, Color: 3},
			{Number: 6, Color: 4},
		},
		{
			{Number: 13, Color: 2},
			{Number: 13, Color: 3},
			{Number: 13, Color: 1},
		},
	}
	if !CanOpenTiles(opened) {
		t.Error("Expected valid opening with total >= 101")
	} else {
		t.Log("PASS Valid Opening")
	}
}

func TestCanOpenTiles_InvalidScore(t *testing.T) {
	opened := [][]Model.Tile{
		{
			{Number: 1, Color: 1},
			{Number: 2, Color: 1},
			{Number: 3, Color: 1},
		},
	}
	if CanOpenTiles(opened) {
		t.Error("Expected invalid opening (score < 101)")
	} else {
		t.Log("PASS Invalid Opening")
	}
}

// Bu ESKI => Gosterge Mavi 4, Mavi 5 = Okey, Joker = Mavi 5
// GECERLI GUNCEL => Gosterge Mavi 11, Mavi 12 = Okey, Joker = Mavi 12
func TestCanOpenTilesWithOkeyAndJoker_Valid(t *testing.T) {
	opened := [][]Model.Tile{
		{
			{Number: 10, Color: ColorEnum.Blue},
			{Number: 13, Color: ColorEnum.Blue},
			{Number: 12, Color: ColorEnum.Blue},
			{Number: 11, Color: ColorEnum.Blue},
			{IsOkey: true},
		},
		{
			{Number: 6, Color: ColorEnum.Yellow},
			{Number: 6, Color: ColorEnum.Red},
			//{IsOkey: true}, // okey, 0 puan
			{Number: 6, Color: ColorEnum.Black},
		},
		{
			{Number: 5, Color: ColorEnum.Red},
			{Number: 5, Color: ColorEnum.Yellow},
			{Number: 5, Color: ColorEnum.Blue}, // 0 puan
		},
		{
			{Number: 12, Color: ColorEnum.Yellow},
			{Number: 12, Color: ColorEnum.Blue, IsJoker: true},
			{Number: 12, Color: ColorEnum.Red},
		},
		{
			{Number: 7, Color: ColorEnum.Red},
			{Number: 8, Color: ColorEnum.Red},
			{Number: 9, Color: ColorEnum.Red},
		},
	}
	if !CanOpenTiles(opened) {
		t.Error("Expected valid opening with total >= 101")
	} else {
		t.Log("PASS Valid Opening")
	}
}
