package Model

import "time"

type AuditLog struct {
	DateTime  time.Time `json:"DateTime"`
	Timestamp time.Time `json:"TimeStamp"` // Mapping'teki "TimeStamp" ile birebir aynı

	OrderID    int64  `json:"OrderID"`
	UserID     int64  `json:"UserID"`
	UserName   string `json:"UserName"`
	ActionType int    `json:"ActionType"`
	ActionName string `json:"ActionName"`
	Message    string `json:"Message"`
	ModuleName string `json:"ModuleName"`

	GameID string `json:"GameID,omitempty"`
	RoomID string `json:"RoomID,omitempty"`

	Tiles *[]Tile `json:"Tiles,omitempty"`

	PenaltyReasonID   *int     `json:"PenaltyReasonID,omitempty"`
	PenaltyReason     *string  `json:"PenaltyReason,omitempty"`
	PenaltyMultiplier *float64 `json:"PenaltyMultiplier,omitempty"`
	PenaltyPoints     *int     `json:"PenaltyPoints,omitempty"`

	HadOkeyTile            *bool `json:"HadOkeyTile,omitempty"`
	OpenedFivePairsButLost *bool `json:"OpenedFivePairsButLost,omitempty"`
	OkeyUsedInFinish       *bool `json:"OkeyUsedInFinish,omitempty"`

	ReconnectDelaySeconds     *float64 `json:"ReconnectDelaySeconds,omitempty"`
	GameDurationSeconds       *float64 `json:"GameDurationSeconds,omitempty"`
	PlayerReactionTimeSeconds *float64 `json:"PlayerReactionTimeSeconds,omitempty"`

	IPAddress string `json:"IPAddress,omitempty"` // Mapping'te "ip" tipi
	Browser   string `json:"Browser,omitempty"`
	Device    string `json:"Device,omitempty"`
	Platform  string `json:"Platform,omitempty"`

	ErrorCode *int                   `json:"ErrorCode,omitempty"`
	ExtraData map[string]interface{} `json:"ExtraData,omitempty"`
}
