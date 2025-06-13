package main

import (
	"fmt"
	"log"
	"okey101/Core"
	"okey101/Elastic"
	"okey101/Model"
	"okey101/Mongo"
	"okey101/Shared"
	"strings"
	"time"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
const playerCount int = 4
const startTileContPerPlayer int = 21
const userIdLength = 12

func main() {
	/*var password, _ = Core.Decrypt("8ii4hYPUQNPziS1PdwhRaevqymWj2eI=", Shared.Config.SECRETKEY)
	fmt.Println(password)
	var elasticUrl, _ = Core.Decrypt("322thoSHDK/x1X4TJVYpMadamuKERhjbwpF3", Shared.Config.SECRETKEY)
	fmt.Println(elasticUrl)*/

	//MongoDB Creation If Not Exist
	/*client, ctx, errMongo := Mongo.MongoOpen()
	if errMongo != nil {
		panic(errMongo)
	}
	defer client.Disconnect(ctx)

	if err := Mongo.EnsureDatabaseAndCollection(client, ctx); err != nil {
		fmt.Println("DB/Collection oluşturulamadı:", err)
	} else {
		fmt.Println("101Okey DB ve Log_Entry collection hazır.")
	}*/
	//-------------------------------
	playerList := make([]Model.Player, playerCount)

	tiles := Core.CreateFullTileSet()
	fmt.Println("Toplam taş:", len(tiles)) // 106 olmalı
	dealUserTiles(playerList, tiles)
	//var player1 = Core.ShowPlayerTiles(&tiles, "Player 1", 22)
	//var player2 = Core.ShowPlayerTiles(&tiles, "Player 2", 21)

	//MongoDB Insert Log----------------
	client, ctx, errMongo := Mongo.MongoOpen()
	if errMongo != nil {
		//log.Fatal("Mongo bağlantı hatası:", errMongo)
		fmt.Println("Mongo bağlantı hatası:", errMongo)
	} else {
		defer client.Disconnect(ctx)
		// LogEntry oluştur
		entry := Mongo.LogEntry{
			DateTime:                  time.Now(),
			TimeStamp:                 time.Now(),
			OrderID:                   1001,
			UserName:                  "player1",
			UserID:                    42,
			ActionType:                1,
			ActionName:                "StartGame",
			Message:                   "Oyun başladı",
			ModuleName:                "GameModule",
			GameID:                    "game123",
			RoomID:                    "room456",
			Tiles:                     Mongo.ConvertCoreTilesToMongoTiles(playerList[0].TileBag), // *[]Tile => []Tile
			PenaltyReasonID:           0,
			PenaltyReason:             "",
			PenaltyMultiplier:         1.0,
			PenaltyPoints:             0,
			HadOkeyTile:               true,
			OpenedFivePairsButLost:    false,
			OkeyUsedInFinish:          false,
			ReconnectDelaySeconds:     0,
			GameDurationSeconds:       360.5,
			PlayerReactionTimeSeconds: 1.75,
			IPAddress:                 "192.168.1.15",
			Browser:                   "Chrome",
			Device:                    "PC",
			Platform:                  "Windows",
			ErrorCode:                 0,
			ExtraData: map[string]interface{}{
				"customField": "customValue",
			},
		}

		id, errMng := Mongo.InsertLogEntry(client, ctx, entry)
		if errMng != nil {
			log.Fatal("LogEntry eklenemedi:", errMng)
		}

		fmt.Println("LogEntry eklendi. ID:", id)
	}
	//MONGODB INSERT LOG FINISH-----------------------------
	//var player2 = Core.ShowPlayerTiles(&tiles, "Player 2:", 21)
	indicatorTile := tiles.GetRandomIndicatorFromTiles()
	fmt.Println("Indicator Tile:")
	fmt.Println(strings.Repeat("-", 30))
	colorName := Core.GetEnumName(Core.ColorEnum, indicatorTile.Color)
	fmt.Printf("ID: %d, %s %d, Joker: %v\n", indicatorTile.ID, colorName, indicatorTile.Number, indicatorTile.IsJoker)
	fmt.Println(strings.Repeat("-", 30))

	fmt.Println("Bag: ")
	for i, tile := range tiles {
		colorName := Core.GetEnumName(Core.ColorEnum, tile.Color)
		//fmt.Printf("%d-) ID: %d, %s %d, Joker: %v\n", i+1, tile.ID, colorName, tile.Number, tile.IsJoker)
		fmt.Printf("%d-) ID: %d, %s %d, Joker: %v Okey: %v\n", i+1, tile.ID, colorName, tile.Number, tile.IsJoker, tile.IsOkey)
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
	//var takenTile = tiles.TakeOneFromBag((*[]Model.Tile)(player1))
	var takenTile = tiles.TakeOneFromBag(&playerList[0].TileBag)
	// Add ElasticLog
	auditLog := Model.AuditLog{
		DateTime:   time.Now(),
		Timestamp:  time.Now(),
		OrderID:    11,   // Oyuna özel bir ID
		UserID:     9876, // Oynayan oyuncunun ID'si
		UserName:   "player1",
		ActionType: int(Core.ActionType.DrawFromMiddle),
		ActionName: Core.GetEnumName(Core.ActionType, Core.ActionType.DrawFromMiddle),
		Message:    fmt.Sprintf("Player %s took a tile from the bag", "player1"),
		ModuleName: "TileBag",

		GameID: "game-abc-123",
		RoomID: "room-xyz-789",

		Tiles: &[]Model.Tile{takenTile}, // Alınan taşın bilgisi loglanır

		IPAddress:                 "192.168.1.100", // Mapping'e göre 'ClientIP' değil, 'IPAddress'
		Browser:                   "Chrome",
		Device:                    "Desktop",
		Platform:                  "Windows",
		PlayerReactionTimeSeconds: Core.FloatPtr(1.2),
		GameDurationSeconds:       Core.FloatPtr(0),
		PenaltyReasonID:           nil,
		PenaltyReason:             nil,
		PenaltyMultiplier:         nil,
		PenaltyPoints:             nil,
		HadOkeyTile:               nil,
		OpenedFivePairsButLost:    nil,
		OkeyUsedInFinish:          nil,

		ErrorCode: nil,
		ExtraData: map[string]interface{}{
			"custom_info": "example",
		},
	}

	err := Elastic.InsertAuditLog(auditLog, Shared.Config.ELASTICAUDITINDEX)
	if err != nil {
		fmt.Println("Audit log insert failed:", err)
	}
	// ElasticLog Bitti--------------------------

	colorName3 := Core.GetEnumName(Core.ColorEnum, takenTile.Color)
	fmt.Printf("Ortadan Cekilen Tas - ID: %d, %s %d, Joker: %v\n", takenTile.ID, colorName3, takenTile.Number, takenTile.IsJoker)

	var dropTile = playerList[0].TileBag[3] //Player1 Bir tas Cantadan cekti ve kendine ekledi
	colorName2 := Core.GetEnumName(Core.ColorEnum, dropTile.Color)
	fmt.Printf("Player1'den Cekilen Tas - ID: %d, %s %d, Joker: %v\n", dropTile.ID, colorName2, dropTile.Number, dropTile.IsJoker)

	Core.DropTileFromTiles((*[]Model.Tile)(&playerList[0].TileBag), dropTile) //Player 1 tas cantadan ceker ve ustten 3.'yu, atar :)
	//fmt.Printf("Player1'den Cekilen Tas - ID: %d, %s %d, Joker: %v\n", dropTile.ID, colorName2, dropTile.Number, dropTile.IsJoker)
	colorNameTaken := Core.GetEnumName(Core.ColorEnum, takenTile.Color)
	fmt.Printf("Player1'den Cekilen Tas - ID: %d, %s %d, Joker: %v\n", takenTile.ID, colorNameTaken, takenTile.Number, takenTile.IsJoker)

	var result = Core.DropTileFromTiles(&playerList[0].TileBag, dropTile) //Player 1 tas cantadan ceker ve ustten 3.'yu, atar :)
	if result {
		fmt.Printf("Player1'den Atilan Tas - ID: %d, %s %d, Joker: %v\n", dropTile.ID, colorName2, dropTile.Number, dropTile.IsJoker)
	}
	// Add ElasticLog for DropTile
	auditLogDrop := Model.AuditLog{
		DateTime:   time.Now(),
		Timestamp:  time.Now(),
		OrderID:    12,   // Yeni unique işlem ID'si
		UserID:     9876, // Oyuncu ID'si
		UserName:   "player1",
		ActionType: int(Core.ActionType.DiscardTile),
		ActionName: Core.GetEnumName(Core.ActionType, Core.ActionType.DiscardTile),
		Message:    fmt.Sprintf("Player %s dropped a tile", "player1"),
		ModuleName: "DropLogic",

		GameID: "game-abc-123",
		RoomID: "room-xyz-789",

		Tiles: &[]Model.Tile{dropTile}, // Atılan taş loglanır

		IPAddress:                 "192.168.1.100",
		Browser:                   "Chrome",
		Device:                    "Desktop",
		Platform:                  "Windows",
		PlayerReactionTimeSeconds: Core.FloatPtr(2.3),
		GameDurationSeconds:       Core.FloatPtr(120.0),

		PenaltyReasonID:        nil,
		PenaltyReason:          nil,
		PenaltyMultiplier:      nil,
		PenaltyPoints:          nil,
		HadOkeyTile:            nil,
		OpenedFivePairsButLost: nil,
		OkeyUsedInFinish:       nil,

		ErrorCode: nil,
		ExtraData: map[string]interface{}{
			"position": 3,
			"source":   "hand",
		},
	}

	err2 := Elastic.InsertAuditLog(auditLogDrop, Shared.Config.ELASTICAUDITINDEX)
	if err2 != nil {
		fmt.Println("Audit log insert failed:", err2)
	}
	// ElasticLog Bitti (DropTile)

	Core.TakeOneFromTable((*[]Model.Tile)(&playerList[1].TileBag), dropTile)                    //Player 2 => Player 1'in 3. elemanini ceker
	Core.DropTileFromTiles((*[]Model.Tile)(&playerList[1].TileBag), (playerList[0].TileBag)[4]) // player 2 ustten 4. elemani atar.
	for i, player := range playerList {
		start := 0
		if i == 0 {
			start = 1
		}
		tb := Core.TileBag(player.TileBag) // []Model.Tile → TileBag
		Core.ShowPlayerTiles(&tb, player.Name, startTileContPerPlayer+start)
	}
	//Core.ShowPlayerTiles(player1, "Player 1", 22)
	//Core.ShowPlayerTiles(player2, "Player 2", 21)
	fmt.Println()
	Core.ShowPlayerTiles(&tiles, "Bag", len(tiles))

}

func dealUserTiles(playerList []Model.Player, tiles Core.TileBag) {
	for i := range playerCount {
		start := 0
		if i == 0 {
			start = 1
		}
		uid, err := Core.GenerateID(userIdLength)
		if err != nil {

		}
		playerList[i] = Model.Player{
			ID:       i,
			Name:     fmt.Sprintf("Player %d", i+1),
			TileBag:  *Core.ShowPlayerTiles(&tiles, fmt.Sprintf("Player %d", i+1), startTileContPerPlayer+start),
			UniqueId: uid,
		}
	}
}
