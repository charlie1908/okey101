package Model

type Tile struct {
	ID      int  `json:"ID"`
	Number  int  `json:"Number"`
	Color   int  `json:"Color"`
	IsJoker bool `json:"IsJoker"`
	IsOkey  bool `json:"IsOkey"`
	IsOpend bool `json:"IsOpend"`
	GroupID *int `json:"GroupID,omitempty"` // sadece açık taşlar için atanır
	X       *int `json:"X,omitempty"`       // UI sıralaması için (isteğe bağlı)
	Y       *int `json:"Y,omitempty"`       // UI grubu için (isteğe bağlı)
}

//IsOpend Property Eklenince Etkilenen Funcionlar
/* - HasAtLeastFivePairs
   - CanOpenTiles
   - CanAddTilesToSet
   - IsValidPair => CanAddPairToPairSets
   - CanThrowingTileBeAddedToOpponentSets
*/
