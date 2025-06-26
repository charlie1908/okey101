package Core

import (
	"fmt"
	"okey101/Model"
	"testing"
	"time"
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

// Test CanOpenTiles And SAVE REDIS USER STATE
// Bu ESKI => Gosterge Mavi 4, Mavi 5 = Okey, Joker = Mavi 5
// GECERLI GUNCEL => Gosterge Mavi 11, Mavi 12 = Okey, Joker = Mavi 12
func TestCanOpenTilesWithOkeyAndJoker_Valid(t *testing.T) {
	roomID := "room123"
	userID := "user456"
	userName := "TestPlayer"

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

		// REDIS TEST SAVE RoomState, PlayerPrivateState, PlayerPublicState
		err := SaveOpenedGroupsToRedis(roomID, userID, userName, opened)
		if err != nil {
			t.Errorf("Failed to save opened groups to Redis: %v", err)
		}
	}
}

func SaveOpenedGroupsToRedis(roomID, userID, userName string, opened [][]*Model.Tile) error {
	//client := GetRedisClient()

	// []*Tile → []Tile (deep copy)
	var openedGroups [][]Model.Tile
	for _, group := range opened {
		var flat []Model.Tile
		for _, tile := range group {
			if tile != nil {
				flat = append(flat, *tile)
			}
		}
		openedGroups = append(openedGroups, flat)
	}
	// 1. PlayerPublicState'i Redis'e yaz. Tasla acildi Player Public bilgisi degisti..
	public := Model.PlayerPublicState{
		UserID:       userID,
		UserName:     userName,
		DiscardTiles: nil,
		OpenedGroups: openedGroups,
		Score:        0,
		IsConnected:  true,
		IsFinished:   false,
		LastDrawTile: nil,
	}

	/*keyPlayerPublic := GeneratePlayerPublicStateRedisKey(roomID, userID)
	err := client.SetKey(keyPlayerPublic, public, 30*time.Minute)
	if err != nil {
		return err
	}*/

	// 2. PlayerPrivateState'i Redis'e yaz => Tas acilinca eldei taslar azaldi..
	private := Model.PlayerPrivateState{
		RoomID:   roomID,
		GameID:   "game789",
		UserID:   userID,
		UserName: userName,
		PlayerTiles: []Model.Tile{
			{Number: 5, Color: ColorEnum.Yellow},
			{Number: 5, Color: ColorEnum.Red},
			{Number: 5, Color: ColorEnum.Blue},
			{Number: 11, Color: ColorEnum.Blue},
			{Number: 12, Color: ColorEnum.Blue, IsOkey: true},
		},
	}

	/*keyPrivate := GeneratePlayerPrivateStateRedisKey(roomID, userID)
	if err := client.SetKey(keyPrivate, private, 30*time.Minute); err != nil {
		return err
	}*/

	//Bu opsiyonel yazildi..
	// 3. RoomState sahte verisi oluşturulup Redis’e yazılıyor
	roomState := Model.RoomState{
		RoomID:    "room123",
		GameID:    "game789",
		Indicator: Model.Tile{Number: 11, Color: ColorEnum.Blue},               // Gösterge taşı
		OkeyTile:  Model.Tile{Number: 12, Color: ColorEnum.Blue, IsOkey: true}, // Okey taşı
		TileBag: []Model.Tile{
			{Number: 1, Color: ColorEnum.Red},
			{Number: 2, Color: ColorEnum.Red},
			{Number: 3, Color: ColorEnum.Red},
			{Number: 4, Color: ColorEnum.Red},
		}, // Sahte taş torbası
		CurrentTurn:   "user456", // TestPlayer’ın sırası
		TurnStartTime: time.Now().UnixMilli(),
		CreatedAt:     time.Now().UnixMilli(),
		GamePhase:     "in_progress",
		WinnerID:      "", // henüz kazanan yok

		Players: []Model.PlayerBasicInfo{
			{UserID: "user456", UserName: "TestPlayer"},
			{UserID: "user123", UserName: "Opponent1"},
			{UserID: "user789", UserName: "Opponent2"},
			{UserID: "user321", UserName: "Opponent3"},
		},
	}

	/*keyRoom := GenerateRoomStateRedisKey(roomID)
	return client.SetKey(keyRoom, roomState, 30*time.Minute)*/
	return SaveGameToRedis(&roomState, []Model.PlayerPrivateState{private}, []Model.PlayerPublicState{public})
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

func TestHasAtLeastFivePairs_Invalid_DoubleJoker(t *testing.T) {
	opened := [][]*Model.Tile{
		{{Number: 8, Color: ColorEnum.Yellow, IsJoker: true}, {Number: 8, Color: ColorEnum.Yellow, IsJoker: true}}, // Sahte okey çifti
		{{Number: 3, Color: ColorEnum.Red}, {Number: 3, Color: ColorEnum.Red}},                                     // Çift 1
		{{Number: 7, Color: ColorEnum.Yellow}, {Number: 7, Color: ColorEnum.Yellow}},                               // Çift 2
		{{Number: 12, Color: ColorEnum.Blue}, {Number: 12, Color: ColorEnum.Blue}},                                 // Çift 3
		{{Number: 5, Color: ColorEnum.Red}, {Number: 5, Color: ColorEnum.Red}},                                     // Çift 4
	}

	if HasAtLeastFivePairs(opened) {
		t.Error("Expected false: two jokers should not form a valid pair")
	} else {
		t.Log("PASS: Double joker pair correctly rejected")
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

	validGroups, remaining, maxScore := SplitTilesByValidGroupsOrRuns(tiles)

	totalSum := sumAllGroupsNumbers(validGroups)
	fmt.Printf("Valid groups toplam sayısı: %d\n MaxScore: %d\n", totalSum, maxScore)

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

	validGroups, remaining, maxScore := SplitTilesByValidGroupsOrRuns(tiles)

	totalSum := sumAllGroupsNumbers(validGroups)
	fmt.Printf("Valid groups toplam sayısı: %d\n MaxScore: %d\n", totalSum, maxScore)

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

	validGroups, remaining, maxScore := SplitTilesByValidGroupsOrRuns(tiles)

	totalSum := sumAllGroupsNumbers(validGroups)
	fmt.Printf("Valid groups toplam sayısı: %d\n MaxScore: %d\n", totalSum, maxScore)

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

	validGroups, remaining, maxScore := SplitTilesByValidGroupsOrRuns(tiles)

	totalSum := sumAllGroupsNumbers(validGroups)
	fmt.Printf("Valid groups toplam sayısı: %d\n MaxScore: %d\n", totalSum, maxScore)

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

	validGroups, remaining, maxScore := SplitTilesByValidGroupsOrRuns(tiles)

	totalSum := sumAllGroupsNumbers(validGroups)
	fmt.Printf("Valid groups toplam sayısı: %d\n MaxScore: %d\n", totalSum, maxScore)

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
	validGroups, remaining, maxScore := SplitTilesByValidGroupsOrRuns(tiles)

	// Grup ve seri kontrolü
	foundGroup := false
	foundSequence := false

	// Kombinasyonların sayı toplamı
	totalSum := sumAllGroupsNumbers(validGroups)
	//t.Logf("🧮 Geçerli grupların toplam sayı değeri: %d", totalSum)
	t.Logf("🧮 Geçerli grupların toplam sayı değeri: %d MaxScore: %d", totalSum, maxScore)

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

func TestGetRemainingInTiles_Valid(t *testing.T) {
	tiles := []*Model.Tile{
		{Color: ColorEnum.Black, Number: 1},
		{Color: ColorEnum.Black, Number: 2},
		{Color: ColorEnum.Black, Number: 3},

		{Color: ColorEnum.Blue, Number: 5},
		{Color: ColorEnum.Blue, Number: 5},

		{Color: ColorEnum.Yellow, Number: 11},
		{Color: ColorEnum.Yellow, Number: 12},
		{Color: ColorEnum.Yellow, Number: 13},

		{Color: ColorEnum.Red, Number: 7},
		{Color: ColorEnum.Red, Number: 8},
		{Color: ColorEnum.Red, Number: 9},

		{Color: ColorEnum.Black, Number: 8},
		{Color: ColorEnum.Black, Number: 8},

		{Color: ColorEnum.Red, Number: 1},
		{Color: ColorEnum.Blue, Number: 1},
		{Color: ColorEnum.Yellow, Number: 1},

		{Color: ColorEnum.Red, Number: 6},
		{Color: ColorEnum.Red, Number: 10},

		{Color: ColorEnum.Blue, Number: 7},
		{Color: ColorEnum.Yellow, Number: 3},
		{Color: ColorEnum.Blue, Number: 11},
	}

	opened := [][]*Model.Tile{
		// Açılmışlar:
		{{Number: 1, Color: ColorEnum.Black}, {Number: 2, Color: ColorEnum.Black}, {Number: 3, Color: ColorEnum.Black}},
		{{Number: 5, Color: ColorEnum.Blue}, {Number: 5, Color: ColorEnum.Blue}},
		{{Number: 11, Color: ColorEnum.Yellow}, {Number: 12, Color: ColorEnum.Yellow}, {Number: 13, Color: ColorEnum.Yellow}},
	}

	remaining := getRemainingInOpenedTiles(tiles, opened)

	if len(remaining) != 13 {
		t.Errorf("Expected 13 remaining tiles, got %d", len(remaining))
	} else {
		t.Log("PASS TestGetRemainingInTiles_Valid - 13 remaining tiles")
	}

	// (İsteğe bağlı) içerik doğrulaması yapılabilir
}

func TestCanOpenTilesWithRemaining_Valid(t *testing.T) {
	tiles := []*Model.Tile{
		// Per: 1 kırmızı, mavi, sarı
		{Number: 1, Color: ColorEnum.Red},
		{Number: 1, Color: ColorEnum.Blue},
		{Number: 1, Color: ColorEnum.Yellow},

		// Seri: 3,4,5 kırmızı
		{Number: 3, Color: ColorEnum.Red},
		{Number: 4, Color: ColorEnum.Red},
		{Number: 5, Color: ColorEnum.Red},

		// Per: 7 sarı, kırmızı, mavi
		{Number: 7, Color: ColorEnum.Yellow},
		{Number: 7, Color: ColorEnum.Red},
		{Number: 7, Color: ColorEnum.Blue},

		// Seri: 8,9,10 siyah
		{Number: 8, Color: ColorEnum.Black},
		{Number: 9, Color: ColorEnum.Black},
		{Number: 10, Color: ColorEnum.Black},

		// Seri: 10,11,12 kırmızı
		{Number: 10, Color: ColorEnum.Red},
		{Number: 11, Color: ColorEnum.Red},
		{Number: 12, Color: ColorEnum.Red},

		// Per: 13 kırmızı, mavi, sarı
		{Number: 13, Color: ColorEnum.Red},
		{Number: 13, Color: ColorEnum.Blue},
		{Number: 13, Color: ColorEnum.Yellow},

		// Kalan taşlar
		{Number: 2, Color: ColorEnum.Yellow},
		{Number: 6, Color: ColorEnum.Blue},
	}

	opened := [][]*Model.Tile{
		/*// Per 1'ler
		{{Number: 1, Color: ColorEnum.Red}, {Number: 1, Color: ColorEnum.Blue}, {Number: 1, Color: ColorEnum.Yellow}},

		// Seri 3,4,5 kırmızı
		{{Number: 3, Color: ColorEnum.Red}, {Number: 4, Color: ColorEnum.Red}, {Number: 5, Color: ColorEnum.Red}},
		*/
		// Per 7'ler
		{{Number: 7, Color: ColorEnum.Yellow}, {Number: 7, Color: ColorEnum.Red}, {Number: 7, Color: ColorEnum.Blue}},

		// Seri 8,9,10 siyah
		{{Number: 8, Color: ColorEnum.Black}, {Number: 9, Color: ColorEnum.Black}, {Number: 10, Color: ColorEnum.Black}},

		// Seri 10,11,12 kırmızı
		{{Number: 10, Color: ColorEnum.Red}, {Number: 11, Color: ColorEnum.Red}, {Number: 12, Color: ColorEnum.Red}},

		// Per 13'ler
		{{Number: 13, Color: ColorEnum.Red}, {Number: 13, Color: ColorEnum.Blue}, {Number: 13, Color: ColorEnum.Yellow}},
	}

	remaining, score, ok := CanOpenTilesWithRemaining(tiles, opened)

	if !ok {
		t.Error("Expected to open tiles with score >= 101, but failed")
	} else {
		t.Logf("PASS CanOpenTilesWithRemaining_Valid: Score = %d, Remaining = %d tiles", score, len(remaining))
		t.Log("Remaining tiles:")
		for _, tile := range remaining {
			t.Logf("Number: %d, Color: %v", tile.Number, GetEnumName(ColorEnum, tile.Color))
		}
	}
}

func Test_CanOpenTilesWithRemainingWithAllGroups_MixedTypesAndLengths(t *testing.T) {
	group1 := []*Model.Tile{
		{Color: ColorEnum.Red, Number: 5},
		{Color: ColorEnum.Blue, Number: 5},
		{Color: ColorEnum.Black, Number: 5},
	}

	group2 := []*Model.Tile{
		{Color: ColorEnum.Yellow, Number: 9},
		{Color: ColorEnum.Red, Number: 9},
		{Color: ColorEnum.Black, Number: 9},
	}

	group3 := []*Model.Tile{
		{Color: ColorEnum.Blue, Number: 1},
		{Color: ColorEnum.Blue, Number: 2},
		{Color: ColorEnum.Blue, Number: 3},
	}

	group4 := []*Model.Tile{
		{Color: ColorEnum.Yellow, Number: 10},
		{Color: ColorEnum.Yellow, Number: 11},
		{Color: ColorEnum.Yellow, Number: 12},
		{Color: ColorEnum.Yellow, Number: 13},
	}

	invalidGroup := []*Model.Tile{
		{Color: ColorEnum.Red, Number: 6},
		{Color: ColorEnum.Blue, Number: 11},
		{Color: ColorEnum.Yellow, Number: 2},
		{Color: ColorEnum.Black, Number: 1},
	}

	invalidGroup2 := []*Model.Tile{
		{Color: ColorEnum.Red, Number: 3},
		{Color: ColorEnum.Blue, Number: 5},
		{Color: ColorEnum.Yellow, Number: 7},
		{Color: ColorEnum.Black, Number: 8},
	}

	allGroups := [][]*Model.Tile{
		group1, group2, group3, group4, invalidGroup, invalidGroup2,
	}

	openedGroups, remainingGroups, score, ok := CanOpenTilesWithRemainingWithAllGroups(allGroups)

	if !ok {
		t.Error("Expected ok == true, got false")
	} else {
		t.Log("PASS ok == true")
	}

	// Açılan grup sayısı 4 olmalı
	if len(openedGroups) != 4 {
		t.Errorf("Expected 4 opened groups, got %d", len(openedGroups))
	} else {
		t.Log("PASS opened groups count == 4")
	}

	// Açılan taş sayısı toplamı 13 olmalı
	totalOpenedTiles := 0
	for _, grp := range openedGroups {
		totalOpenedTiles += len(grp)
	}
	if totalOpenedTiles != 13 {
		t.Errorf("Expected total opened tiles 13, got %d", totalOpenedTiles)
	} else {
		t.Log("PASS total opened tiles == 13")
	}

	// Remaining grup sayısı 1 olmalı
	if len(remainingGroups) != 1 {
		t.Errorf("Expected 1 remaining group, got %d", len(remainingGroups))
	} else {
		t.Log("PASS remaining groups count == 1")
	}

	// Remaining grupundeki taş sayısı 8 olmalı
	if len(remainingGroups[0]) != 8 {
		t.Errorf("Expected remaining group size 8, got %d", len(remainingGroups[0]))
	} else {
		t.Log("PASS remaining group size == 8")
	}

	// Skor > 0 olmalı
	if score <= 0 {
		t.Errorf("Expected score > 0, got %d", score)
	} else {
		t.Logf("PASS score > 0 (%d)", score)
	}

	// remaining içinde invalidGroup'un ilk taşı var mı kontrolü
	found := false
	for _, tile := range remainingGroups[0] {
		if tile == invalidGroup[0] {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected remaining group to contain tiles from invalidGroup")
	} else {
		t.Log("PASS remaining group contains invalidGroup tiles")
	}
}

func TestRedisLoadGameStateForPlayer(t *testing.T) {
	client := GetRedisClient()
	roomID := "room_test_123"
	userID := "user_test_1"

	// 1. Dummy veriler hazırlanır
	room := Model.RoomState{
		RoomID:        roomID,
		GameID:        "game_001",
		Indicator:     Model.Tile{Number: 5, Color: ColorEnum.Blue},
		OkeyTile:      Model.Tile{Number: 6, Color: ColorEnum.Blue, IsOkey: true},
		TileBag:       []Model.Tile{{Number: 1, Color: ColorEnum.Red}},
		CurrentTurn:   userID,
		TurnStartTime: time.Now().UnixMilli(),
		CreatedAt:     time.Now().UnixMilli(),
		GamePhase:     "in_progress",
		WinnerID:      "",
		Players: []Model.PlayerBasicInfo{
			{UserID: userID, UserName: "TestUser1"},
			{UserID: "user_test_2", UserName: "TestUser2"},
		},
	}

	private := Model.PlayerPrivateState{
		RoomID:      roomID,
		GameID:      "game_001",
		UserID:      userID,
		UserName:    "TestUser1",
		PlayerTiles: []Model.Tile{{Number: 3, Color: ColorEnum.Red}, {Number: 4, Color: ColorEnum.Red}},
	}

	public1 := Model.PlayerPublicState{
		UserID:       userID,
		UserName:     "TestUser1",
		DiscardTiles: []Model.Tile{{Number: 9, Color: ColorEnum.Yellow}},
		OpenedGroups: nil,
		Score:        0,
		IsConnected:  true,
		IsFinished:   false,
		LastDrawTile: nil,
	}

	public2 := Model.PlayerPublicState{
		UserID:       "user_test_2",
		UserName:     "TestUser2",
		DiscardTiles: []Model.Tile{{Number: 8, Color: ColorEnum.Black}},
		OpenedGroups: nil,
		Score:        0,
		IsConnected:  true,
		IsFinished:   false,
		LastDrawTile: nil,
	}

	// 2. Redis'e veriler yazılır
	if err := client.SetKey(GenerateRoomStateRedisKey(roomID), room, 30*time.Minute); err != nil {
		t.Fatal("Room state Redis'e yazılamadı:", err)
	}
	if err := client.SetKey(GeneratePlayerPrivateStateRedisKey(roomID, userID), private, 30*time.Minute); err != nil {
		t.Fatal("Private state Redis'e yazılamadı:", err)
	}
	if err := client.SetKey(GeneratePlayerPublicStateRedisKey(roomID, userID), public1, 30*time.Minute); err != nil {
		t.Fatal("Public state (1) Redis'e yazılamadı:", err)
	}
	if err := client.SetKey(GeneratePlayerPublicStateRedisKey(roomID, "user_test_2"), public2, 30*time.Minute); err != nil {
		t.Fatal("Public state (2) Redis'e yazılamadı:", err)
	}

	//Wait 5 seconds
	//time.Sleep(5 * time.Second)

	// 3. Fonksiyon test edilir
	loadedRoom, loadedPrivate, loadedPublics, err := LoadGameForPlayer(roomID, userID)
	if err != nil {
		t.Fatal("LoadGameForPlayer başarısız:", err)
	}

	// 4. Veriler doğrulanır
	if loadedRoom.RoomID != roomID {
		t.Error("RoomID eşleşmiyor")
	}
	if loadedPrivate.UserID != userID || loadedPrivate.UserName != "TestUser1" {
		t.Error("Private state hatalı")
	}
	if len(loadedPublics) != 2 {
		t.Errorf("Beklenen public state sayısı 2, gelen: %d", len(loadedPublics))
	}

	found := false
	for _, pub := range loadedPublics {
		if pub.UserID == "user_test_2" && pub.UserName == "TestUser2" {
			found = true
		}
	}
	if !found {
		t.Error("Public state içinde user_test_2 bulunamadı")
	}

	//Player 1 Elapsed Time
	elapsed := (time.Now().UnixMilli() - room.TurnStartTime) / 1000
	if elapsed > 30 {
		// oyuncunun süresi doldu
		fmt.Println("Player1 Elapsed Time Already Finished: ", elapsed)
	} else {
		fmt.Println("Player1 Elapsed Time: ", elapsed)
	}

	t.Log("PASS: LoadGameForPlayer doğru şekilde çalıştı")
}

func TestGetEffectiveNumber_OkeyWithSequenceEndAt13(t *testing.T) {
	// Grup: 11 (Mavi), 12 (Mavi), 13 (Mavi) + 1 Okey taşı
	group := []*Model.Tile{
		{Number: 11, Color: ColorEnum.Blue},
		{Number: 12, Color: ColorEnum.Blue},
		{Number: 13, Color: ColorEnum.Blue},
		{Number: 5, Color: ColorEnum.Red, IsOkey: true}, // Okey taşı
	}

	okeyTile := group[len(group)-1] // Okey taşı referansı

	// getEffectiveNumber fonksiyonunu çağırıyoruz
	effectiveNumber := getEffectiveNumber(okeyTile, group)

	expected := 10 // Sonraki sayı 14 değil, baştan 1 eksik sayı olarak 10 dönmeli

	if effectiveNumber != expected {
		t.Errorf("Beklenen Okey değeri %d, bulundu: %d", expected, effectiveNumber)
	} else {
		t.Logf("✅ Okey değeri doğru hesaplandı: %d", effectiveNumber)
	}

	// Okey taşını da dahil ederek grubu efektif değere göre sırala
	sortedGroup := sortGroupByEffectiveNumber(group)

	t.Log("Sıralı grup:")
	for _, tile := range sortedGroup {
		t.Logf("Number=%d, Color=%v, IsOkey=%v", tile.Number, tile.Color, tile.IsOkey)
	}

	// Ek olarak, isSequence kontrolü
	nonOkeys := filterNonOkeys(group)
	okeyCount := countOkeys(group)

	if !isSequence(nonOkeys, okeyCount) {
		t.Error("❌ Bu grup seri olarak tanımlanmalıydı ama isSequence false döndü.")
	} else {
		t.Log("✅ Seri doğru tanımlandı.")
	}
}
