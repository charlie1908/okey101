package Model

type Tile struct {
	ID      int  `json:"ID"`
	Number  int  `json:"Number"`
	Color   int  `json:"Color"`
	IsJoker bool `json:"IsJoker"`
	IsOkey  bool `json:"IsOkey"`
	IsOpend bool `json:"IsOpend"`
}

//IsOpend Property Eklenince Etkilenen Funcionlar
/* - HasAtLeastFivePairs
   - CanOpenTiles
   - CanAddTilesToSet
   - IsValidPair => CanAddPairToPairSets
   - CanThrowingTileBeAddedToOpponentSets
*/
