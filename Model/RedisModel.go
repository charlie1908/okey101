package Model

type RoomState struct {
	RoomID      string
	GameID      string
	Indicator   Tile
	OkeyTile    Tile
	TileBag     []Tile // Kalan taşlar
	CurrentTurn string // Sıra kimde
	CreatedAt   int64
}
type PlayerState struct {
	RoomID       string
	GameID       string
	UserID       string
	UserName     string
	PlayerTiles  []Tile
	DiscardTiles []Tile
	IsConnected  bool
}
