package Model

type Tile struct {
	ID      int  // Her taşın eşsiz kimliği (örneğin: 0-105 arası, 106 sahte okey)
	Number  int  // 1–13 arası
	Color   int  // Renk (Red, Yellow, Blue, Black)
	IsJoker bool // Sahte okey mi
	IsOkey  bool // Gerçek okey mi (göstergeye göre belirlenir)
}
