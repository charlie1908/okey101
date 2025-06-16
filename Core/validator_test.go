package Core

import (
	"fmt"
	"okey101/Model"
	"testing"
)

// Test Group gecerli olmali.
// Gecersiz yani false oldugu durumda Test yanlistir.
func TestValidGroup(t *testing.T) {
	tiles := []*Model.Tile{
		{Number: 5, Color: ColorEnum.Red},
		{Number: 5, Color: ColorEnum.Blue},
		{Number: 5, Color: ColorEnum.Yellow},
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
	tiles := []*Model.Tile{
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
	tiles := []*Model.Tile{
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
	tiles := []*Model.Tile{
		{Number: 12, Color: ColorEnum.Blue, IsOkey: true},
		{Number: 12, Color: ColorEnum.Blue, IsOkey: true},
		{Number: 7, Color: 1},
		{Number: 4, Color: 1},
	}
	if !IsValidGroupOrRun(tiles) {
		t.Error("Expected valid run with joker, got invalid")
	} else {
		t.Log("PASS Valid Run with Okey")
	}
}

func TestCanOpenTiles_Valid(t *testing.T) {
	opened := [][]*Model.Tile{
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
	opened := [][]*Model.Tile{
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
	opened := [][]*Model.Tile{
		{
			{Number: 10, Color: ColorEnum.Blue},
			{Number: 13, Color: ColorEnum.Blue},
			{Number: 12, Color: ColorEnum.Blue},
			{Number: 11, Color: ColorEnum.Blue},
			//{IsOkey: true},
			{Number: 12, Color: ColorEnum.Blue, IsOkey: true},
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

// GECERLI GUNCEL => Kirmizi 5, Kirmizi 6 = Okey, Joker = Kirmizi 6
func TestCanOpenTilesWithOkeyAndJokerSameLineGroup_Valid(t *testing.T) {
	opened := [][]*Model.Tile{
		{
			{Number: 10, Color: ColorEnum.Blue},
			{Number: 13, Color: ColorEnum.Blue},
			{Number: 12, Color: ColorEnum.Blue},
			{Number: 11, Color: ColorEnum.Blue},
		},
		{
			{Number: 6, Color: ColorEnum.Yellow},
			{Number: 12, Color: ColorEnum.Blue, IsOkey: true}, // okey
			{Number: 6, Color: ColorEnum.Red, IsJoker: true},  //Aslinda Burda bir hata var Joker Mavi 12 bu hata yakalanmali. => *Note
		},
		{
			{Number: 5, Color: ColorEnum.Red},
			{Number: 5, Color: ColorEnum.Yellow},
			{Number: 5, Color: ColorEnum.Blue},
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

// GECERLI GUNCEL => Blue 11, Blue 12 = Okey, Joker = Blue 12
func TestCanOpenTilesWithOkeyAndJokerSameLineSquence_Valid(t *testing.T) {
	opened := [][]*Model.Tile{
		{
			{Number: 10, Color: ColorEnum.Blue},
			{Number: 11, Color: ColorEnum.Blue, IsJoker: true},
			{Number: 11, Color: ColorEnum.Red, IsOkey: true},
		},
		{
			{Number: 11, Color: ColorEnum.Yellow},
			{Number: 11, Color: ColorEnum.Red, IsOkey: true}, // Okey
			{Number: 11, Color: ColorEnum.Red, IsJoker: true},
		},
		{
			{Number: 5, Color: ColorEnum.Red},
			{Number: 5, Color: ColorEnum.Yellow},
			{Number: 5, Color: ColorEnum.Blue}, // 0 puan
		},
		{
			{Number: 9, Color: ColorEnum.Yellow},
			//{Number: 9, Color: ColorEnum.Blue, IsJoker: true},
			{Number: 9, Color: ColorEnum.Blue},
			{Number: 9, Color: ColorEnum.Red},
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

//************************************

func TestHasAtLeastFivePairs_Valid(t *testing.T) {
	opened := [][]*Model.Tile{
		{{Number: 8, Color: ColorEnum.Yellow, IsOkey: true}, {Number: 8, Color: ColorEnum.Yellow, IsOkey: true}},
		{{Number: 3, Color: ColorEnum.Red}, {Number: 3, Color: ColorEnum.Red}},       // Çift 1
		{{Number: 7, Color: ColorEnum.Yellow}, {Number: 7, Color: ColorEnum.Yellow}}, // Çift 2
		{{Number: 12, Color: ColorEnum.Blue}, {Number: 12, Color: ColorEnum.Blue}},   // Çift 3
		{{Number: 5, Color: ColorEnum.Red}, {Number: 5, Color: ColorEnum.Red}},       // Çift 4
		//{{Number: 9, Color: ColorEnum.Black}, {IsOkey: true}},                        // Okey ile çift 5 (9 siyah + okey)
	}

	if !HasAtLeastFivePairs(opened) {
		t.Error("Expected to have at least five valid pairs")
	} else {
		t.Log("PASS HasAtLeastFivePairs with valid 5 pairs")
	}
}

func TestHasAtLeastFivePairsIsOpened_Invalid(t *testing.T) {
	opened := [][]*Model.Tile{
		{{Number: 3, Color: ColorEnum.Red}, {Number: 3, Color: ColorEnum.Red, IsOpend: true}}, // Biri onceden acilmis
		{{Number: 7, Color: ColorEnum.Yellow}, {Number: 7, Color: ColorEnum.Yellow}},
		{{Number: 12, Color: ColorEnum.Blue}, {Number: 12, Color: ColorEnum.Blue}},
		{{Number: 5, Color: ColorEnum.Red}, {Number: 5, Color: ColorEnum.Red}},
		{{Number: 9, Color: ColorEnum.Black}, {Number: 9, Color: ColorEnum.Black}},
	}

	if HasAtLeastFivePairs(opened) {
		t.Error("Expected to NOT have five valid pairs")
	} else {
		t.Log("PASS HasAtLeastFivePairs with invalid pairs")
	}
}

func TestHasAtLeastFivePairs_Invalid(t *testing.T) {
	opened := [][]*Model.Tile{
		{{Number: 3, Color: ColorEnum.Red}, {Number: 3, Color: ColorEnum.Yellow}},    // Renk farklı, çift değil
		{{Number: 7, Color: ColorEnum.Yellow}, {Number: 8, Color: ColorEnum.Yellow}}, // Sayı farklı
		{{Number: 12, Color: ColorEnum.Blue}, {Number: 12, Color: ColorEnum.Black}},  // Renk farklı
		{{Number: 5, Color: ColorEnum.Red}},                                          // Tek taş, çift değil
		{{Number: 9, Color: ColorEnum.Black}, {Number: 9, Color: ColorEnum.Blue}},    // Renk farklı
	}

	if HasAtLeastFivePairs(opened) {
		t.Error("Expected to NOT have five valid pairs")
	} else {
		t.Log("PASS HasAtLeastFivePairs with invalid pairs")
	}
}

func TestHasLeastWithSixPairs_InValid(t *testing.T) {
	opened := [][]*Model.Tile{
		{
			{Number: 13, Color: ColorEnum.Red},
			{Number: 13, Color: ColorEnum.Red},
		},
		{
			//{Number: 12, Color: ColorEnum.Blue},
			{IsOkey: true, Number: 5, Color: ColorEnum.Red},
			{Number: 12, Color: ColorEnum.Blue},
		},
		{
			{Number: 11, Color: ColorEnum.Yellow},
			{Number: 11, Color: ColorEnum.Yellow},
		},
		{
			//{Number: 10, Color: ColorEnum.Black},
			{IsOkey: true, Number: 5, Color: ColorEnum.Red},
			{IsJoker: true, Number: 4, Color: ColorEnum.Red},
			//{Number: 10, Color: ColorEnum.Black},
		},
		{
			{Number: 9, Color: ColorEnum.Red},
			{Number: 9, Color: ColorEnum.Red},
		},
		{
			//Evet 5 pair tutuyor ama 6. pair hatali..
			{Number: 3, Color: ColorEnum.Blue},
			{Number: 3, Color: ColorEnum.Black},
		},
	}

	if HasAtLeastFivePairs(opened) {
		t.Error("Expected valid five-pair opening")
	} else {
		t.Log("PASS InValid Six Pairs")
	}
}

// *****Bu hatali gecmemeli
func TestHasLeastWithFivePairs_InvalidLessThanFivePairs(t *testing.T) {
	opened := [][]*Model.Tile{
		{
			{Number: 5, Color: ColorEnum.Red},
			{Number: 5, Color: ColorEnum.Red},
		},
		{
			{Number: 6, Color: ColorEnum.Blue},
			{Number: 6, Color: ColorEnum.Blue},
		},
		{
			{Number: 7, Color: ColorEnum.Yellow},
			{Number: 7, Color: ColorEnum.Yellow},
		},
	}

	if HasAtLeastFivePairs(opened) {
		t.Error("Expected invalid opening (less than five pairs)")
	} else {
		t.Log("PASS Invalid - Less Than Five Pairs")
	}
}

//************************************

func TestCanAddTilesToSet_ValidRunAddition(t *testing.T) {
	gid := 5
	set := []*Model.Tile{
		{Number: 4, Color: ColorEnum.Blue, IsOpend: true, GroupID: &gid},
		{Number: 5, Color: ColorEnum.Blue, IsOpend: true, GroupID: &gid},
		{Number: 6, Color: ColorEnum.Blue, IsOpend: true, GroupID: &gid},
	}
	newTile := Model.Tile{Number: 7, Color: ColorEnum.Blue}

	if !CanAddTilesToSet(set, &newTile) {
		t.Error("Expected to successfully add tile to run")
	} else {
		t.Log("PASS: Added tile to run successfully")
	}
}

func TestCanAddJokerToSet_ValidRunAddition(t *testing.T) {
	gid := 6
	set := []*Model.Tile{
		/*{Number: 11, Color: ColorEnum.Blue},
		{Number: 12, Color: ColorEnum.Blue},
		{Number: 13, Color: ColorEnum.Blue},*/
		{Number: 5, Color: ColorEnum.Blue, IsOpend: true, GroupID: &gid},
		{Number: 6, Color: ColorEnum.Blue, IsOpend: true, GroupID: &gid},
		{Number: 7, Color: ColorEnum.Blue, IsOpend: true, GroupID: &gid},
	}
	newTile := Model.Tile{Number: 4, Color: ColorEnum.Blue, IsOkey: true}

	if !CanAddTilesToSet(set, &newTile) {
		t.Error("Expected to successfully add tile to run")
	} else {
		//0 ise basa basa bir sayi ise ve arada bosluk yok ise son rakkam 13'de degil ise sona koy
		//var OkeyTileValue = CalculateTileScore(&newTile, 0, set, true)
		var OkeyTileValue = CalculateTileScore(&newTile, 2, set, true)
		t.Log("Okey Tile Value: ", OkeyTileValue)
		t.Log("PASS: Added tile to run successfully")
	}
}

func TestCanAddTilesToSet_InvalidAddToGroup(t *testing.T) {
	gid := 7
	set := []*Model.Tile{
		{Number: 8, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid},
		{Number: 8, Color: ColorEnum.Yellow, IsOpend: true, GroupID: &gid},
		{Number: 8, Color: ColorEnum.Blue, IsOpend: true, GroupID: &gid},
	}
	newTile := Model.Tile{Number: 8, Color: ColorEnum.Red}

	if CanAddTilesToSet(set, &newTile) {
		t.Error("Expected failure when adding tile to a pair set (less than 3 tiles)")
	} else {
		t.Log("PASS: Cannot add tile to pair set")
	}
}

func TestCanAddTilesToSetIsOpened_InvalidAddToGroup(t *testing.T) {
	gid := 8
	set := []*Model.Tile{
		{Number: 8, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid},
		{Number: 8, Color: ColorEnum.Yellow, IsOpend: true, GroupID: &gid},
		{Number: 8, Color: ColorEnum.Blue, IsOpend: true, GroupID: &gid},
	}
	newTile := Model.Tile{Number: 8, Color: ColorEnum.Black, IsOpend: true}

	if CanAddTilesToSet(set, &newTile) {
		t.Error("Expected failure when adding tile to a pair set (less than 3 tiles)")
	} else {
		t.Log("PASS: Cannot add tile to pair set")
	}
}

//************************************

// Cifte Git kullaniciya Tas Isleme
func TestCanAddPairToPairSets_ValidPairWithEnoughPairs(t *testing.T) {
	// Geçerli çift
	newPair := []*Model.Tile{
		{Number: 8, Color: ColorEnum.Red, IsOpend: false},
		{Number: 8, Color: ColorEnum.Red, IsOpend: false},
	}

	gid := 9
	gid2 := 10
	gid3 := 11
	gid4 := 12
	gid5 := 13
	// En az 5 çift açılmış setler
	pairSets := [][]*Model.Tile{
		{
			{Number: 1, Color: ColorEnum.Blue, IsOpend: true, GroupID: &gid},
			{Number: 1, Color: ColorEnum.Blue, IsOpend: true, GroupID: &gid},
		},
		{
			{Number: 2, Color: ColorEnum.Yellow, IsOpend: true, GroupID: &gid2},
			{Number: 2, Color: ColorEnum.Yellow, IsOpend: true, GroupID: &gid2},
		},
		{
			{Number: 3, Color: ColorEnum.Black, IsOpend: true, GroupID: &gid3},
			{Number: 3, Color: ColorEnum.Black, IsOpend: true, GroupID: &gid3},
		},
		{
			{Number: 4, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid4},
			{Number: 4, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid4},
		},
		{
			{Number: 5, Color: ColorEnum.Blue, IsOpend: true, GroupID: &gid5},
			{Number: 5, Color: ColorEnum.Blue, IsOpend: true, GroupID: &gid5},
		},
	}

	if !CanAddPairToPairSets(newPair, pairSets) {
		t.Error("Bekleniyor: geçerli çift ve en az 5 çift açılmış set olduğunda true dönmeli")
	} else {
		t.Log("PASS: Geçerli çift ve yeterli çift set ile çift eklenebilir")
	}
}

func TestCanAddPairToPairSets_InvalidPairOrNotEnoughPairs(t *testing.T) {
	// Geçersiz çift (farklı renkler, aynı sayı değil)
	newPair := []*Model.Tile{
		{Number: 8, Color: ColorEnum.Red, IsOpend: false},
		{Number: 9, Color: ColorEnum.Red, IsOpend: false},
		{Number: 10, Color: ColorEnum.Red, IsOpend: false},
	}
	/*newPair := []*Model.Tile{
		{Number: 8, Color: ColorEnum.Red, IsOpend: true},
		{Number: 8, Color: ColorEnum.Red, IsOpend: false},
	}*/

	gid := 9
	gid2 := 10
	gid3 := 11
	gid4 := 12
	gid5 := 13
	// 4 çift açılmış setler (yeterli değil)
	pairSets := [][]*Model.Tile{
		{
			{Number: 1, Color: ColorEnum.Blue, IsOpend: true, GroupID: &gid},
			{Number: 1, Color: ColorEnum.Blue, IsOpend: true, GroupID: &gid},
		},
		{
			{Number: 2, Color: ColorEnum.Yellow, IsOpend: true, GroupID: &gid2},
			{Number: 2, Color: ColorEnum.Yellow, IsOpend: true, GroupID: &gid2},
		},
		{
			{Number: 3, Color: ColorEnum.Black, IsOpend: true, GroupID: &gid3},
			{Number: 3, Color: ColorEnum.Black, IsOpend: true, GroupID: &gid3},
		},
		{
			{Number: 4, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid4},
			{Number: 4, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid4},
		},
		{
			{Number: 5, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid5},
			{Number: 5, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid5},
		},
	}

	if CanAddPairToPairSets(newPair, pairSets) {
		t.Error("Bekleniyor: geçersiz çift veya yeterli çift yoksa false dönmeli")
	} else {
		t.Log("PASS: Geçersiz çift veya yetersiz çift set ile çift eklenemez")
	}
}

//************************************

func TestCanThrowingTileBeAddedToOpponentSets_ValidAddition_TotalOver101(t *testing.T) {
	gid := 9
	gid2 := 10
	gid3 := 11
	opponentSets := [][]*Model.Tile{
		{
			{Number: 13, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid},
			{Number: 13, Color: ColorEnum.Blue, IsOpend: true, GroupID: &gid},
			{Number: 13, Color: ColorEnum.Yellow, IsOpend: true, GroupID: &gid}, // Grup: Aynı sayı, farklı renk
		},
		{
			{Number: 11, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid2},
			{Number: 12, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid2},
			{Number: 13, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid2}, // Seri: Artan sayılar, aynı renk
		},
		{
			{Number: 8, Color: ColorEnum.Black, IsOpend: true, GroupID: &gid3},
			{Number: 9, Color: ColorEnum.Black, IsOpend: true, GroupID: &gid3},
			{Number: 10, Color: ColorEnum.Black, IsOpend: true, GroupID: &gid3},
			{Number: 11, Color: ColorEnum.Black, IsOpend: true, GroupID: &gid3},
		},
	}
	newTile := Model.Tile{Number: 13, Color: ColorEnum.Black, IsOpend: false} // Eksik rengi tamamlayan taş

	//Atilan Tas Eklenebildigi Group'un GroupID degerini alir!. Boylece nereye eklenebilecegi anlasilabilir ve ceza puani verilir.
	if !CanThrowingTileBeAddedToOpponentSets(&newTile, opponentSets) {
		t.Error("Expected tile to be added to one of the opponent's sets")
	} else {
		t.Log("PASS: Tile can be added to an opponent's set")
	}
}

func TestCanThrowingTileBeAddedToOpponentSets_InvalidAddition_TotalOver101(t *testing.T) {
	gid := 9
	gid2 := 10
	gid3 := 11
	opponentSets := [][]*Model.Tile{
		{
			{Number: 13, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid},
			{Number: 12, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid},
			{Number: 11, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid}, // Seri → 36
		},
		{
			{Number: 13, Color: ColorEnum.Blue, IsOpend: true, GroupID: &gid2},
			{Number: 13, Color: ColorEnum.Yellow, IsOpend: true, GroupID: &gid2},
			{Number: 13, Color: ColorEnum.Black, IsOpend: true, GroupID: &gid2}, // Grup → 39
		},
		{
			{Number: 10, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid3},
			{Number: 9, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid3},
			{Number: 8, Color: ColorEnum.Red, IsOpend: true, GroupID: &gid3}, // Seri → 27
		},
	}
	newTile := Model.Tile{Number: 7, Color: ColorEnum.Yellow, IsOpend: false} // Hiçbir sete uymaz

	//Bu Test'de atilan tas hicbir group'a uymadigi icin IsOpened: false ve GroupID'si nil dir!
	if CanThrowingTileBeAddedToOpponentSets(&newTile, opponentSets) {
		t.Error("Expected tile NOT to be added to any of the opponent's sets")
	} else {
		t.Log("PASS: Tile cannot be added to any of the opponent's sets")
	}
}

func TestValidGroupAndRunWithNoJokers(t *testing.T) {
	tiles := []*Model.Tile{
		{Number: 9, Color: ColorEnum.Yellow}, // remaining
		{Number: 4, Color: ColorEnum.Black},  // group
		{Number: 7, Color: ColorEnum.Black},  // remaining
		{Number: 10, Color: ColorEnum.Blue},  // run 2
		{Number: 6, Color: ColorEnum.Yellow}, // run 1
		{Number: 3, Color: ColorEnum.Black},  // remaining
		{Number: 8, Color: ColorEnum.Red},    // group
		{Number: 1, Color: ColorEnum.Yellow}, // remaining
		{Number: 11, Color: ColorEnum.Blue},  // run 2
		{Number: 7, Color: ColorEnum.Yellow}, // run 1
		{Number: 5, Color: ColorEnum.Yellow}, // run 1
		{Number: 4, Color: ColorEnum.Blue},   // group
		{Number: 5, Color: ColorEnum.Red},    // remaining
		{Number: 12, Color: ColorEnum.Blue},  // run 2
		{Number: 6, Color: ColorEnum.Black},  // remaining
		{Number: 2, Color: ColorEnum.Blue},   // remaining
		{Number: 1, Color: ColorEnum.Red},    // remaining
		{Number: 9, Color: ColorEnum.Blue},   // remaining
		{Number: 7, Color: ColorEnum.Red},    // remaining
		{Number: 4, Color: ColorEnum.Yellow}, // run 1
		{Number: 13, Color: ColorEnum.Red},   // remaining
	}

	validGroups, remaining := SplitTilesByValidGroupsOrRuns(tiles)

	totalSum := sumAllGroupsNumbers(validGroups)
	fmt.Printf("Valid groups toplam sayısı: %d\n", totalSum)

	foundGroup := false
	foundRun := false

	for i, group := range validGroups {
		if isGroup(group, 0) {
			foundGroup = true
			t.Logf("Geçerli Grup #%d:", i+1)
		} else if isSequence(group, 0) {
			foundRun = true
			t.Logf("Geçerli Seri #%d:", i+1)
		}
		for _, tile := range group {
			t.Logf(" - Taş: Number=%d, Color=%v", tile.Number, GetEnumName(ColorEnum, tile.Color))
		}
	}

	if !foundGroup {
		t.Errorf("❌ Geçerli bir grup bulunamadı (aynı sayı, farklı renk).")
	}
	if !foundRun {
		t.Errorf("❌ Geçerli bir sıra bulunamadı (aynı renk, ardışık).")
	}
	if len(remaining) == 0 {
		t.Errorf("❌ Kalan taş kalmadı, oysa en az 1 tane geçersiz olmalıydı.")
	} else {
		t.Logf("✅ %d geçersiz taş remaining içinde kaldı.", len(remaining))
	}
}

func TestValidGroupAndRunWithOkeysAndJokers(t *testing.T) {
	tiles := []*Model.Tile{
		{Number: 7, Color: ColorEnum.Blue},                // Group/Run aday
		{Number: 5, Color: ColorEnum.Blue},                // Run aday
		{Number: 7, Color: ColorEnum.Red},                 // Group/Run aday
		{Number: 9, Color: ColorEnum.Yellow},              // Remaining
		{Number: 6, Color: ColorEnum.Blue},                // Run aday
		{Number: 7, Color: ColorEnum.Black},               // Group/Run aday
		{Number: 11, Color: ColorEnum.Red},                // Run aday
		{Number: 4, Color: ColorEnum.Blue},                // Run aday
		{Number: 10, Color: ColorEnum.Red},                // Remaining
		{Number: 10, Color: ColorEnum.Blue, IsOkey: true}, // Okey 1
		{Number: 3, Color: ColorEnum.Black},               // Remaining
		{Number: 10, Color: ColorEnum.Blue, IsOkey: true}, // Okey 2
		{Number: 8, Color: ColorEnum.Blue, IsJoker: true}, // Joker
		{Number: 1, Color: ColorEnum.Yellow},              // Remaining
	}

	validGroups, remaining := SplitTilesByValidGroupsOrRuns(tiles)

	totalSum := sumAllGroupsNumbers(validGroups)
	fmt.Printf("Valid groups toplam sayısı: %d\n", totalSum)

	foundGroup := false
	foundRun := false

	for i, group := range validGroups {
		if isGroup(filterNonOkeys(group), countOkeys(group)) {
			foundGroup = true
			t.Logf("Geçerli Grup #%d:", i+1)
		} else if isSequence(filterNonOkeys(group), countOkeys(group)) {
			foundRun = true
			t.Logf("Geçerli Seri #%d:", i+1)
		}
		for _, tile := range group {
			t.Logf(" - Taş: Number=%d, Color=%v, IsOkey=%v, IsJoker=%v", tile.Number, GetEnumName(ColorEnum, tile.Color), tile.IsOkey, tile.IsJoker)
		}
	}

	if !foundGroup {
		t.Errorf("❌ Geçerli bir grup bulunamadı (aynı sayı, farklı renk, Okey/Joker dahil).")
	}
	if !foundRun {
		t.Errorf("❌ Geçerli bir sıra bulunamadı (aynı renk, ardışık, Okey/Joker dahil).")
	}
	if len(remaining) == 0 {
		t.Errorf("❌ Kalan taş kalmadı, oysa en az 1 tane geçersiz olmalıydı.")
	} else {
		t.Logf("✅ %d geçersiz taş remaining içinde kaldı.", len(remaining))
	}
}

func TestSplitTilesWithOkeyAndJokerInGroupsAndSequences(t *testing.T) {
	tiles := []*Model.Tile{
		{Number: 7, Color: ColorEnum.Red},
		{Number: 5, Color: ColorEnum.Blue},
		{Number: 10, Color: ColorEnum.Red},
		{Number: 10, Color: ColorEnum.Blue},
		{Number: 3, Color: ColorEnum.Black},
		{Number: 6, Color: ColorEnum.Blue},
		{Number: 3, Color: ColorEnum.Yellow},
		{Number: 8, Color: ColorEnum.Blue},
		{Number: 10, Color: ColorEnum.Yellow},
		{Number: 3, Color: ColorEnum.Red},
		{Number: 6, Color: ColorEnum.Red},
		{Number: 8, Color: ColorEnum.Red},
		{Number: 5, Color: ColorEnum.Yellow},
		{Number: 7, Color: ColorEnum.Blue},
		{Number: 1, Color: ColorEnum.Yellow},
		{Number: 12, Color: ColorEnum.Blue},
		{Number: 9, Color: ColorEnum.Yellow},
		{Number: 10, Color: ColorEnum.Blue, IsOkey: true}, // Okey
		{Number: 11, Color: ColorEnum.Blue},
		{Number: 6, Color: ColorEnum.Yellow},
		{Number: 10, Color: ColorEnum.Blue, IsJoker: true},
	}

	validGroups, remaining := SplitTilesByValidGroupsOrRuns(tiles)

	totalSum := sumAllGroupsNumbers(validGroups)
	fmt.Printf("Valid groups toplam sayısı: %d\n", totalSum)

	foundGroup := false
	foundSequence := false

	for _, group := range validGroups {
		if isGroup(filterNonOkeys(group), countOkeys(group)) {
			foundGroup = true
			t.Logf("✅ Geçerli Grup Bulundu:")
			for _, tile := range group {
				t.Logf("  - Taş: Number=%d, Color=%v, IsOkey=%v, IsJoker=%v", tile.Number, GetEnumName(ColorEnum, tile.Color), tile.IsOkey, tile.IsJoker)
			}
		} else if isSequence(filterNonOkeys(group), countOkeys(group)) {
			foundSequence = true
			t.Logf("✅ Geçerli Seri Bulundu:")
			for _, tile := range group {
				t.Logf("  - Taş: Number=%d, Color=%v, IsOkey=%v, IsJoker=%v", tile.Number, GetEnumName(ColorEnum, tile.Color), tile.IsOkey, tile.IsJoker)
			}
		}
	}

	if !foundGroup {
		t.Error("❌ Geçerli grup bulunamadı.")
	}
	if !foundSequence {
		t.Error("❌ Geçerli seri bulunamadı.")
	}
	if len(remaining) == 0 {
		t.Error("❌ Kalan taş olmamalıydı.")
	} else {
		t.Logf("ℹ️ %d taş geçersiz kaldı.", len(remaining))
	}
}

func TestSplitTilesWithOkeyAndJokerInOneSquenceAndNoGroup(t *testing.T) {
	tiles := []*Model.Tile{
		{Number: 5, Color: ColorEnum.Blue},
		{Number: 9, Color: ColorEnum.Red},
		{Number: 6, Color: ColorEnum.Blue},
		{Number: 1, Color: ColorEnum.Black},
		{Number: 3, Color: ColorEnum.Red},
		{Number: 9, Color: ColorEnum.Black},
		{Number: 2, Color: ColorEnum.Yellow},
		{Number: 8, Color: ColorEnum.Blue, IsJoker: true}, // Joker (4,5,6 serisi için 4 yerine)
		{Number: 9, Color: ColorEnum.Blue, IsOkey: true},  // Okey (9'lu grup için)
		{Number: 4, Color: ColorEnum.Yellow},
		{Number: 12, Color: ColorEnum.Red},
		{Number: 7, Color: ColorEnum.Red},
		{Number: 11, Color: ColorEnum.Blue},
		{Number: 10, Color: ColorEnum.Black},
		{Number: 13, Color: ColorEnum.Yellow},
		{Number: 1, Color: ColorEnum.Red},
		{Number: 3, Color: ColorEnum.Black},
		{Number: 7, Color: ColorEnum.Black},
		{Number: 8, Color: ColorEnum.Yellow},
		{Number: 2, Color: ColorEnum.Blue},
		{Number: 6, Color: ColorEnum.Yellow},
	}

	validGroups, remaining := SplitTilesByValidGroupsOrRuns(tiles)

	totalSum := sumAllGroupsNumbers(validGroups)
	fmt.Printf("Valid groups toplam sayısı: %d\n", totalSum)

	foundGroup := false
	foundSequence := false

	for _, group := range validGroups {
		if isGroup(filterNonOkeys(group), countOkeys(group)) {
			foundGroup = true
			t.Logf("✅ Geçerli Grup Bulundu:")
			for _, tile := range group {
				t.Logf("  - Taş: Number=%d, Color=%v, IsOkey=%v, IsJoker=%v", tile.Number, GetEnumName(ColorEnum, tile.Color), tile.IsOkey, tile.IsJoker)
			}
		} else if isSequence(filterNonOkeys(group), countOkeys(group)) {
			foundSequence = true
			t.Logf("✅ Geçerli Seri Bulundu:")
			for _, tile := range group {
				t.Logf("  - Taş: Number=%d, Color=%v, IsOkey=%v, IsJoker=%v", tile.Number, GetEnumName(ColorEnum, tile.Color), tile.IsOkey, tile.IsJoker)
			}
		}
	}

	/*if !foundGroup {
		t.Error("❌ Geçerli grup bulunamadı.")
	}*/
	if foundGroup {
		t.Error("❌ Geçerli grup bulundu!.")
	}
	if !foundSequence {
		t.Error("❌ Geçerli seri bulunamadı.")
	}
	if len(remaining) == 0 {
		t.Error("❌ Kalan taş olmamalıydı.")
	} else {
		t.Logf("ℹ️ %d taş geçersiz kaldı.", len(remaining))
	}
}

func TestSplitTilesWithOkeyAndJokerInOneSquenceAndOneGroup2(t *testing.T) {
	tiles := []*Model.Tile{
		{Number: 9, Color: ColorEnum.Yellow},
		{Number: 5, Color: ColorEnum.Yellow},
		{Number: 8, Color: ColorEnum.Blue, IsJoker: true}, // Joker (normal taş gibi)
		{Number: 9, Color: ColorEnum.Red},
		{Number: 6, Color: ColorEnum.Yellow},
		{Number: 9, Color: ColorEnum.Black},
		{Number: 4, Color: ColorEnum.Yellow},
		{Number: 9, Color: ColorEnum.Blue, IsOkey: true}, // Okey
		{Number: 7, Color: ColorEnum.Red},
		{Number: 6, Color: ColorEnum.Blue},
		{Number: 3, Color: ColorEnum.Red},
		{Number: 1, Color: ColorEnum.Black},
		{Number: 2, Color: ColorEnum.Blue},
		{Number: 11, Color: ColorEnum.Blue},
		{Number: 3, Color: ColorEnum.Black},
		{Number: 10, Color: ColorEnum.Black},
		{Number: 13, Color: ColorEnum.Yellow},
		{Number: 1, Color: ColorEnum.Red},
		{Number: 7, Color: ColorEnum.Black},
		{Number: 8, Color: ColorEnum.Yellow},
		{Number: 2, Color: ColorEnum.Yellow},
	}

	validGroups, remaining := SplitTilesByValidGroupsOrRuns(tiles)

	totalSum := sumAllGroupsNumbers(validGroups)
	fmt.Printf("Valid groups toplam sayısı: %d\n", totalSum)

	foundGroup := false
	foundSequence := false

	for _, group := range validGroups {
		if isGroup(filterNonOkeys(group), countOkeys(group)) {
			foundGroup = true
			t.Logf("✅ Geçerli Grup Bulundu:")
			for _, tile := range group {
				t.Logf("  - Taş: Number=%d, Color=%v, IsOkey=%v, IsJoker=%v", tile.Number, GetEnumName(ColorEnum, tile.Color), tile.IsOkey, tile.IsJoker)
			}
		} else if isSequence(filterNonOkeys(group), countOkeys(group)) {
			foundSequence = true
			t.Logf("✅ Geçerli Seri Bulundu:")
			for _, tile := range group {
				t.Logf("  - Taş: Number=%d, Color=%v, IsOkey=%v, IsJoker=%v", tile.Number, GetEnumName(ColorEnum, tile.Color), tile.IsOkey, tile.IsJoker)
			}
		}
	}

	if !foundGroup {
		t.Error("❌ Geçerli grup bulunamadı.")
	}
	if !foundSequence {
		t.Error("❌ Geçerli seri bulunamadı.")
	}
	if len(remaining) == 0 {
		t.Error("❌ Kalan taş olmamalıydı.")
	} else {
		t.Logf("ℹ️ %d taş geçersiz kaldı.", len(remaining))
	}
}

func TestSplitTilesWithOkeyAndJokerInOneSequenceAndNoGroup(t *testing.T) {
	tiles := []*Model.Tile{
		// Seri (4-5-6-7 Blue)
		{Number: 5, Color: ColorEnum.Blue},
		{Number: 4, Color: ColorEnum.Blue},
		{Number: 6, Color: ColorEnum.Blue},
		{Number: 7, Color: ColorEnum.Blue},

		// Grup (9'luk): 9 Red, 9 Black, 9 Joker, 9 Okey
		{Number: 9, Color: ColorEnum.Red},
		{Number: 9, Color: ColorEnum.Black},
		{Number: 2, Color: ColorEnum.Red, IsJoker: true},   // Joker (9 yerine)
		{Number: 9, Color: ColorEnum.Yellow, IsOkey: true}, // Okey (9 yerine)

		// Diğer taşlar (geçersiz)
		{Number: 1, Color: ColorEnum.Yellow},
		{Number: 2, Color: ColorEnum.Blue},
		{Number: 11, Color: ColorEnum.Black},
		{Number: 13, Color: ColorEnum.Red},
		{Number: 3, Color: ColorEnum.Black},
		{Number: 10, Color: ColorEnum.Yellow},
		{Number: 6, Color: ColorEnum.Red},
		{Number: 7, Color: ColorEnum.Black},
		{Number: 8, Color: ColorEnum.Yellow},
		{Number: 12, Color: ColorEnum.Blue},
		{Number: 1, Color: ColorEnum.Red},
		{Number: 3, Color: ColorEnum.Yellow},
		{Number: 5, Color: ColorEnum.Black},
	}

	// Geçerli kombinasyonları bul
	validGroups, remaining := SplitTilesByValidGroupsOrRuns(tiles)

	// Grup ve seri kontrolü
	foundGroup := false
	foundSequence := false

	// Kombinasyonların sayı toplamı
	totalSum := sumAllGroupsNumbers(validGroups)
	t.Logf("🧮 Geçerli grupların toplam sayı değeri: %d", totalSum)

	// Her kombinasyonu incele
	for _, group := range validGroups {
		nonOkeys := filterNonOkeys(group)
		okeyCount := countOkeys(group)

		if isGroup(nonOkeys, okeyCount) {
			foundGroup = true
			t.Logf("✅ Geçerli Grup Bulundu:")
		} else if isSequence(nonOkeys, okeyCount) {
			foundSequence = true
			t.Logf("✅ Geçerli Seri Bulundu:")
		} else {
			t.Errorf("❌ Ne grup ne de seri: %+v", group)
			continue
		}

		// Taş detaylarını yazdır
		for _, tile := range group {
			t.Logf("   - Taş: Number=%d, Color=%v, IsOkey=%v, IsJoker=%v",
				tile.Number, GetEnumName(ColorEnum, tile.Color), tile.IsOkey, tile.IsJoker)
		}
	}

	// Doğrulamalar
	if foundGroup {
		t.Error("❌ Hic bir grup beklenmiyordu.")
	}
	if !foundSequence {
		t.Error("❌ En az bir geçerli seri bekleniyordu.")
	}
	if len(remaining) == 0 {
		t.Error("❌ Kalan taş olmalıydı.")
	} else {
		t.Logf("ℹ️ %d taş geçersiz kaldı.", len(remaining))
	}
}

func TestSplitTilesByValidPairs_OkeyDifferentFromPairs(t *testing.T) {
	tiles := []*Model.Tile{
		{Number: 7, Color: ColorEnum.Yellow},
		{Number: 3, Color: ColorEnum.Red, IsOkey: true},
		{Number: 5, Color: ColorEnum.Blue},
		{Number: 1, Color: ColorEnum.Red},
		{Number: 9, Color: ColorEnum.Black},
		{Number: 2, Color: ColorEnum.Red},
		{Number: 7, Color: ColorEnum.Yellow},
		{Number: 11, Color: ColorEnum.Blue},
		{Number: 3, Color: ColorEnum.Blue},
		{Number: 5, Color: ColorEnum.Blue, IsJoker: true},
		{Number: 4, Color: ColorEnum.Black},
		{Number: 6, Color: ColorEnum.Yellow},
		{Number: 4, Color: ColorEnum.Yellow},
		{Number: 6, Color: ColorEnum.Yellow},
		{Number: 8, Color: ColorEnum.Black},
		{Number: 4, Color: ColorEnum.Blue},
		{Number: 10, Color: ColorEnum.Red},
		{Number: 3, Color: ColorEnum.Black},
		{Number: 12, Color: ColorEnum.Yellow},
		{Number: 13, Color: ColorEnum.Black},
		{Number: 2, Color: ColorEnum.Blue},
	}
	pairs, remaining := SplitTilesByValidPairs(tiles)
	if len(pairs) != 4 {
		t.Errorf("Beklenen 4 çift, bulunan: %d", len(pairs))
	}
	for i, pair := range pairs {
		if len(pair) != 2 {
			t.Errorf("Pair %d içinde 2 taş olmalı ama %d taş var", i+1, len(pair))
		}
	}
	if len(remaining) == 0 {
		t.Errorf("Kalan taş olmalıydı ama hiç kalmamış")
	}
	for i, pair := range pairs {
		t.Logf("Pair %d:", i+1)
		for _, tile := range pair {
			t.Logf("  - Taş: Number=%d, Color=%v, IsOkey=%v, IsJoker=%v", tile.Number, GetEnumName(ColorEnum, tile.Color), tile.IsOkey, tile.IsJoker)
		}
	}
	t.Logf("Kalan taş sayısı: %d", len(remaining))
	for _, tile := range remaining {
		t.Logf("  - Taş: Number=%d, Color=%v, IsOkey=%v, IsJoker=%v", tile.Number, GetEnumName(ColorEnum, tile.Color), tile.IsOkey, tile.IsJoker)
	}
}

func TestSplitTilesByValidPairs_WithJokerNoOkey(t *testing.T) {
	tiles := []*Model.Tile{
		{Number: 2, Color: ColorEnum.Yellow},
		{Number: 7, Color: ColorEnum.Red},
		{Number: 5, Color: ColorEnum.Blue, IsJoker: true}, // Joker
		{Number: 13, Color: ColorEnum.Black},
		{Number: 11, Color: ColorEnum.Blue},
		{Number: 7, Color: ColorEnum.Yellow},
		{Number: 9, Color: ColorEnum.Yellow},
		{Number: 10, Color: ColorEnum.Black},
		{Number: 6, Color: ColorEnum.Black},
		{Number: 7, Color: ColorEnum.Blue},
		{Number: 8, Color: ColorEnum.Red},
		{Number: 3, Color: ColorEnum.Yellow},
		{Number: 12, Color: ColorEnum.Red},
		{Number: 2, Color: ColorEnum.Blue},
		{Number: 5, Color: ColorEnum.Yellow},
		{Number: 4, Color: ColorEnum.Black},
		{Number: 9, Color: ColorEnum.Yellow},
		{Number: 10, Color: ColorEnum.Blue},
		{Number: 5, Color: ColorEnum.Red},
		{Number: 4, Color: ColorEnum.Black},
		{Number: 12, Color: ColorEnum.Red},
	}

	pairs, remaining := SplitTilesByValidPairs(tiles)

	expectedPairCount := 3
	if len(pairs) != expectedPairCount {
		t.Errorf("❌ Beklenen %d çift, bulunan: %d", expectedPairCount, len(pairs))
	}

	for i, pair := range pairs {
		if len(pair) != 2 {
			t.Errorf("❌ Pair %d içinde 2 taş olmalı ama %d taş var", i+1, len(pair))
		} else {
			t.Logf("✅ Pair %d:", i+1)
			for _, tile := range pair {
				t.Logf("  - Taş: Number=%d, Color=%v, IsJoker=%v",
					tile.Number, GetEnumName(ColorEnum, tile.Color), tile.IsJoker)
			}
		}
	}

	t.Logf("Kalan taş sayısı: %d", len(remaining))
	for _, tile := range remaining {
		t.Logf("  - Taş: Number=%d, Color=%v, IsJoker=%v",
			tile.Number, GetEnumName(ColorEnum, tile.Color), tile.IsJoker)
	}
}
