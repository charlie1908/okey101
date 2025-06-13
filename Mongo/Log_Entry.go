package Mongo

import "time"

type Tile struct {
	ID      int  `bson:"ID"`
	Number  int  `bson:"Number"`
	Color   int  `bson:"Color"`
	IsJoker bool `bson:"IsJoker"`
	IsOkey  bool `bson:"IsOkey"`
}

type LogEntry struct {
	DateTime                  time.Time              `bson:"DateTime"`
	TimeStamp                 time.Time              `bson:"TimeStamp"`
	OrderID                   int64                  `bson:"OrderID"`
	UserName                  string                 `bson:"UserName"`
	UserID                    int64                  `bson:"UserID"`
	ActionType                int                    `bson:"ActionType"`
	ActionName                string                 `bson:"ActionName"`
	Message                   string                 `bson:"Message"`
	ModuleName                string                 `bson:"ModuleName"`
	GameID                    string                 `bson:"GameID"`
	RoomID                    string                 `bson:"RoomID"`
	Tiles                     []Tile                 `bson:"Tiles"`
	PenaltyReasonID           int                    `bson:"PenaltyReasonID"`
	PenaltyReason             string                 `bson:"PenaltyReason"`
	PenaltyMultiplier         float64                `bson:"PenaltyMultiplier"`
	PenaltyPoints             int                    `bson:"PenaltyPoints"`
	HadOkeyTile               bool                   `bson:"HadOkeyTile"`
	OpenedFivePairsButLost    bool                   `bson:"OpenedFivePairsButLost"`
	OkeyUsedInFinish          bool                   `bson:"OkeyUsedInFinish"`
	ReconnectDelaySeconds     float64                `bson:"ReconnectDelaySeconds"`
	GameDurationSeconds       float64                `bson:"GameDurationSeconds"`
	PlayerReactionTimeSeconds float64                `bson:"PlayerReactionTimeSeconds"`
	IPAddress                 string                 `bson:"IPAddress"`
	Browser                   string                 `bson:"Browser"`
	Device                    string                 `bson:"Device"`
	Platform                  string                 `bson:"Platform"`
	ErrorCode                 int                    `bson:"ErrorCode"`
	ExtraData                 map[string]interface{} `bson:"ExtraData"`
}
