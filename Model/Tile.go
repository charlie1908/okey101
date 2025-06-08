package Model

type Tile struct {
	ID      int  `json:"ID"`
	Number  int  `json:"Number"`
	Color   int  `json:"Color"`
	IsJoker bool `json:"IsJoker"`
	IsOkey  bool `json:"IsOkey"`
}
