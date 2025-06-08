package Core

import "reflect"

var ColorEnum = newcolorEnum()
var ActionType = newActionType()
var PenaltyType = newPenaltyType()

func newcolorEnum() *color {
	return &color{
		Red:    1,
		Yellow: 2,
		Blue:   3,
		Black:  4,
		None:   5,
	}
}

type color struct {
	Red    int
	Yellow int
	Blue   int
	Black  int
	None   int
}

func newActionType() *actionType {
	return &actionType{
		DistributedTiles:   0,
		JoinRoom:           1,
		LeaveRoom:          2,
		RejoinRoom:         3,
		KickedRoom:         4,
		Timeout:            5,
		StartGame:          6,
		DiscardTile:        7,
		DrawFromMiddle:     8,
		DrawFromDiscard:    9,
		OpenSet:            10,
		MergeSet:           11,
		FinishGame:         12,
		FinishGameWithOkey: 13,
	}
}

type actionType struct {
	DistributedTiles   int
	JoinRoom           int
	LeaveRoom          int
	RejoinRoom         int
	KickedRoom         int
	Timeout            int
	StartGame          int
	DiscardTile        int
	DrawFromMiddle     int
	DrawFromDiscard    int
	OpenSet            int
	MergeSet           int
	FinishGame         int
	FinishGameWithOkey int
}

func newPenaltyType() *penaltyType {
	return &penaltyType{
		IllegalDiscard:     1,
		OkeyPenalty:        2,
		OpenSetFailed:      3,
		FivePairsButNotWin: 4,
		DoubleDiscard:      5,
		TimeoutPenalty:     6,
	}
}

type penaltyType struct {
	IllegalDiscard     int
	OkeyPenalty        int
	OpenSetFailed      int
	FivePairsButNotWin int
	DoubleDiscard      int
	TimeoutPenalty     int
}

func GetEnumName(enum interface{}, value int) string {
	v := reflect.ValueOf(enum).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if int(v.Field(i).Int()) == value {
			return t.Field(i).Name
		}
	}
	return "Unknown"
}
