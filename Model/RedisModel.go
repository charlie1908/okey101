package Model

type RoomState struct {
	RoomID        string
	GameID        string
	Indicator     Tile
	OkeyTile      Tile
	TileBag       []Tile // Masada kalan taşlar
	CurrentTurn   string // Şu an kimin sırası (UserID) => Belki bunu tutmayiz cunku her seferinde RoomSate'in de guncellenmesi gerekiyor.
	TurnStartTime int64  // Sıra başlama zamanı => TurnStartTime: time.Now().UnixMilli()
	CreatedAt     int64  //CreatedAt: time.Now().UnixMilli()

	GamePhase string //"waiting", "in_progress", "finished" => Core.GamePhase.GamePhaseInProgress
	WinnerID  string //Oyunu kazanan varsa(UserID) Ne zaman ? => GamePass == "finished" durumunda.
	Players   []PlayerBasicInfo
}

type PlayerBasicInfo struct {
	UserID   string
	UserName string
}

type PlayerPublicState struct {
	UserID       string
	UserName     string
	DiscardTiles []Tile   // Masaya atılan taşlar
	OpenedGroups [][]Tile // Açılan setler
	Score        int64
	IsConnected  bool
	IsFinished   bool

	LastDrawTile   *Tile  // Son çekilen taş
	LastDrawSource string // "tilebag", "discard", vs. => Core.GamePhase.LastDrawSource
}

type PlayerPrivateState struct {
	RoomID      string
	GameID      string
	UserID      string
	UserName    string
	PlayerTiles []Tile // Oyuncunun elindeki taşlar
}
