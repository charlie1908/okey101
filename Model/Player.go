package Model

type Player struct {
	ID       int    // Oyuncu No
	UniqueId string // Oyuncu Id si
	Name     string // Oyuncu adı
	TileBag  []Tile // Tile Bag
}
