package Core

import (
	"okey101/Model"
	"sort"
)

// Kendisine gelen tas dizilimlerinin Ayni Sayi Farkjli Renk (Group) veya Sirali Ayni Renk (RUN) olma durumuna bakma.
func IsValidGroupOrRun(tiles []*Model.Tile) bool {
	if len(tiles) < 3 {
		return false
	}

	okeyCount := countOkeys(tiles)
	nonOkeyTiles := filterNonOkeys(tiles)

	return isGroup(nonOkeyTiles, okeyCount) || isSequence(nonOkeyTiles, okeyCount)
}

// Kirmizi 5, Mavi 5, Sari 5
func isGroup(tiles []*Model.Tile, okeyCount int) bool {
	if len(tiles) == 0 {
		return false
	}

	number := tiles[0].Number
	colors := make(map[int]bool)

	for _, tile := range tiles {
		//Ayni sayi degeri degil ise ya da bu renk map'de zaten var ise(tum taslarin rengi farkli olmali) false donulur.
		if tile.Number != number || colors[tile.Color] {
			return false
		}
		colors[tile.Color] = true
	}

	return len(tiles)+okeyCount >= 3
}

// Siyah 3, Siyah 4, Siyah 5
func isSequence(tiles []*Model.Tile, okeyCount int) bool {
	if len(tiles) == 0 {
		return false
	}

	//Ayni renkte degiller ise zaten seri olamaz
	if !allSameColor(tiles) {
		return false
	}

	//Kullanici o seriyi sirali vermese bile siralar dogru mu diye siralayip bakariz..
	//***Buradaki siralama orjinal gelen taslarin sirasini degistimes. Sadece 	nonOkeyTiles := filterNonOkeys(tiles)'den gelen yeni slice tipinin yerini degistirir..
	sort.Slice(tiles, func(i, j int) bool {
		return tiles[i].Number < tiles[j].Number
	})

	//Arada bosluklar var ise gerekli Okey sayisina bakilir..
	neededOkeys := calculateNeededOkeysForRun(tiles, okeyCount)

	//Ihtiyac duyulan OkeyCount eldekinden fazla ise veya toplam serideki tas sayisi 3'den kucuk ise false firlatilir.
	return neededOkeys <= okeyCount && len(tiles)+okeyCount >= 3
}

func calculateNeededOkeysForRun(tiles []*Model.Tile, maxOkeys int) int {
	neededOkeys := 0
	for i := 1; i < len(tiles); i++ {
		diff := tiles[i].Number - tiles[i-1].Number
		switch diff {
		case 1:
			continue //Arada bosluk yok ise
		case 2: //Arada bosluk 1 ise
			neededOkeys++ //Arada bosluk 1 ise
		case 3: //Arada bosluk 2 ise
			neededOkeys += 2
		default:
			return maxOkeys + 1 // Geçersiz yapar
		}
	}
	return neededOkeys
}

func countOkeys(tiles []*Model.Tile) int {
	count := 0
	for _, tile := range tiles {
		if tile.IsOkey {
			count++
		}
	}
	return count
}

func filterNonOkeys(tiles []*Model.Tile) []*Model.Tile {
	result := make([]*Model.Tile, 0, len(tiles))
	for _, tile := range tiles {
		if !tile.IsOkey {
			result = append(result, tile)
		}
	}
	return result
}

func allSameColor(tiles []*Model.Tile) bool {
	if len(tiles) == 0 {
		return true
	}
	color := tiles[0].Color
	for _, tile := range tiles {
		if tile.Color != color {
			return false
		}
	}
	return true
}

//func CalculateTileScore(tile Model.Tile) int {
//if tile.IsOkey || tile.IsJoker {
//	return 0
//}
//return tile.Number
//}

// SQUENCE(ROW)'DA OKEY'I DOGRU YERE KOYMAK ZORUNDA. SIRALI OLMALI! [ 11,OKEY,13 ] GIBI
/*func CalculateTileScore(tile Model.Tile, index int, tiles []Model.Tile, isSequence bool) int {
	if !tile.IsOkey {
		return tile.Number
	}

	//Bundan sonrasi Oeky Tasinin Degerinin hesaplanmasidir.
	if !isSequence {
		//Bir siralama yok tum taslar ayni number degerinde olmali. Hepsinin rengi tabi ki farkli..
		// Group: Okey veya Joker taşları grubun diğer taşları ile aynı number olur.
		// Diğer taşlardan birinin number'ını alabiliriz. Ilk okey veya joker olmayan tasin degeri alinir.
		for _, t := range tiles {
			if !t.IsOkey {
				return t.Number
			}
		}
		// Eğer tüm taşlar okey ise, varsayılan 0 döner.
		return 0
	}

	// Sequence (run) için
	// Aynı renkten artan sıra olduğu için, en küçük number'dan başlayarak index ile artırılır.
	// Öncelikle Okey olmayan taşlardan en küçük number bulunur:
	minNumber := 14 //En yuksek deger alinarak en kucuk okey ve joker olmayan tas bulunur.
	for _, t := range tiles {
		if !t.IsOkey {
			if t.Number < minNumber {
				minNumber = t.Number
			}
		}
	}
	if minNumber == 14 {
		// Tüm taşlar okey ise 0 dönebiliriz.
		return 0
	}

	// Şimdi, index konumuna göre number hesapla:
	// Dizide index 0 ise minNumber, index 1 ise minNumber+1 vs.
	number := minNumber + index
	if number > 13 {
		number = number - 13
	}

	return number
}*/

/*func CalculateTileScore(tile Model.Tile, index int, tiles []Model.Tile, isSequence bool) int {
	if !tile.IsOkey {
		return tile.Number
	}

	if !isSequence {
		// Grup: Diğer taşların number'ı ile aynı olmalı
		//Bir siralama yok tum taslar ayni number degerinde olmali. Hepsinin rengi tabi ki farkli..
		//Group: Okey veya Joker taşları grubun diğer taşları ile aynı number olur.
		//Diğer taşlardan birinin number'ını alabiliriz. Ilk okey veya joker olmayan tasin degeri alinir.
		for _, t := range tiles {
			if !t.IsOkey {
				return t.Number
			}
		}
		return 0
	}
	// Sequence (run) için
	// Aynı renkten artan sıra olduğu için, en küçük number'dan başlayarak index ile artırılır.
	// Öncelikle Okey olmayan taşlardan en küçük number bulunur

	// Sequence (sıra): Okey'in doldurduğu boşluğu bulmalıyız
	nonOkeys := filterNonOkeys(tiles)
	if len(nonOkeys) == 0 {
		return 0
	}

	//Okey olmayan taslar siralanir.
	sort.Slice(nonOkeys, func(i, j int) bool {
		return nonOkeys[i].Number < nonOkeys[j].Number
	})

	minNumber := nonOkeys[0].Number // Ilk deger zaten en kucuk degerdir.
	expectedNumber := minNumber

	for i := 0; i < len(tiles); i++ {
		if i == index {
			return expectedNumber
		}

		if !tiles[i].IsOkey {
			if tiles[i].Number != expectedNumber {
				// Arada Okey kullanılmış gibi davran
				expectedNumber = tiles[i].Number
			}
		}

		expectedNumber++
		if expectedNumber > 13 {
			expectedNumber = 1
		}
	}

	return 0 // Fallback, olmaması gerekir
}*/

func CalculateTileScore(tile *Model.Tile, index int, tiles []*Model.Tile, isSequence bool) int {
	if !tile.IsOkey {
		return tile.Number
	}

	if !isSequence {
		// Grup: Diğer taşların number'ı ile aynı olmalı
		//Bir siralama yok tum taslar ayni number degerinde olmali. Hepsinin rengi tabi ki farkli..
		//Group: Okey veya Joker taşları grubun diğer taşları ile aynı number olur.
		//Diğer taşlardan birinin number'ını alabiliriz. Ilk okey veya joker olmayan tasin degeri alinir.
		for _, t := range tiles {
			if !t.IsOkey {
				return t.Number
			}
		}
		//return 0
		//Butun taslar OK ise kendi degerini alir..
		return tile.Number
	}
	// Sequence (run) için
	// Aynı renkten artan sıra olduğu için, en küçük number'dan başlayarak index ile artırılır.

	// Sequence durumu: Okey'in doğru değerini hesapla

	// Okey olmayan taşları bul
	nonOkeys := []Model.Tile{}
	for _, t := range tiles {
		if !t.IsOkey {
			nonOkeys = append(nonOkeys, *t)
		}
	}

	if len(nonOkeys) == 0 {
		return 1 // Sadece okey varsa 1 dönebiliriz
	}

	// Okey olmayan taşları sirala
	sort.Slice(nonOkeys, func(i, j int) bool {
		return nonOkeys[i].Number < nonOkeys[j].Number
	})

	//Okey olmayan esas taslari numbers'a ata..
	numbers := []int{}
	for _, t := range nonOkeys {
		numbers = append(numbers, t.Number)
	}

	// Aradaki boşluğu bul (fark > 1 olan yer) Sirali gitmeyen degeri bul ve bir sonraki olamsi gereken degeri don. Sonraki - Onceki > 1
	for i := 0; i < len(numbers)-1; i++ {
		if numbers[i+1]-numbers[i] > 1 {
			return numbers[i] + 1
		}
	}

	// Boşluk yoksa:
	if numbers[len(numbers)-1] < 13 { //Sona Koy 13'den buyuk olmuyorsa
		//return numbers[len(numbers)-1] + 1
		if index == 0 {
			//Basa Koy
			return numbers[0] - 1
		} else {
			//Degil ise sona koy
			return numbers[len(numbers)-1] + 1
		}
	}

	if numbers[0] > 1 { // 1'den buyuk ise ve arada bosluk yok se basa koy.
		return numbers[0] - 1
	}

	return 1
}

// Acilan seri uygun ve toplamlari 101 ve uzeri ise..
func CanOpenTiles(opened [][]*Model.Tile) bool {
	totalScore := 0
	for _, group := range opened {

		//Iclerinden Acilan Var Ise Hata Verir..
		if HasOpenTail(group...) {
			return false
		}

		//Once gecerliligi kontrol et!
		if !IsValidGroupOrRun(group) {
			return false
		}

		isSeq := isSequence(filterNonOkeys(group), countOkeys(group))
		for index, tile := range group {

			totalScore += CalculateTileScore(tile, index, group, isSeq)
		}
	}
	var result = totalScore >= 101
	//Acilan Tum Taslar Opened olur!
	if result {
		SetOpentiles(opened)
	}
	//---------
	return result
}

func CanOpenTilesWithRemaining(tiles []*Model.Tile, opened [][]*Model.Tile) (remaining []*Model.Tile, score int, error bool) {
	totalScore := 0
	var remainList []*Model.Tile
	for _, group := range opened {

		//Iclerinden Acilan Var Ise Hata Verir..
		if HasOpenTail(group...) {
			return remainList, 0, false
		}

		//Once gecerliligi kontrol et!
		if !IsValidGroupOrRun(group) {
			return remainList, 0, false
		}

		isSeq := isSequence(filterNonOkeys(group), countOkeys(group))
		for index, tile := range group {

			totalScore += CalculateTileScore(tile, index, group, isSeq)
		}
	}
	var result = totalScore >= 101
	//---------
	//Acilan Tum Taslar belirlendikten sonra geri kalanlar tanimlanir!
	if result {
		remainList = getRemainingInOpenedTiles(tiles, opened)
	}
	//Bayramin istegi ile kaldirildi.
	//Acilan Tum Taslar Opened olur!
	/*if result {
		SetOpentiles(opened)
	}*/
	//---------
	return remainList, totalScore, result
}

func getRemainingInOpenedTiles(tiles []*Model.Tile, opened [][]*Model.Tile) (remaining []*Model.Tile) {
	type pairKey struct {
		Color  int
		Number int
	}

	used := make(map[pairKey][]*Model.Tile)
	var remainList []*Model.Tile

	// Açılmış taşları key'e göre gruplandır
	for _, group := range opened {
		for _, tile := range group {
			key := pairKey{Color: tile.Color, Number: tile.Number}
			used[key] = append(used[key], tile)
		}
	}

	// Elimizdeki taşlardan açılmayanları bul
	for _, tile := range tiles {
		key := pairKey{Color: tile.Color, Number: tile.Number}
		if usedList, ok := used[key]; ok && len(usedList) > 0 {
			// Bu taş kullanıldı, bu taşı tüket ve listeden cikar
			used[key] = usedList[1:]
		} else {
			// Bu taş hiç kullanılmamış
			remainList = append(remainList, tile)
		}
	}
	return remainList
}

// Acilan Tas gurubunun IsOpened'ini true olarak ata.
func SetOpentiles(opened [][]*Model.Tile) {
	for _, group := range opened {
		//101'i gecen elde acilan tum array grouplara ayri ayri unique groupID tanimlanir.
		var groupID = Game.GenerateGroupID()
		for _, tile := range group {
			tile.IsOpend = true
			tile.GroupID = &groupID
		}
	}
}

// Acilan taslarda IsOpened = true olarak ata.
func SetOpenPairtiles(setTiles []*Model.Tile, grpID ...int) {
	//Pair olarak acilan taslara Global ayni siradaki GroupID atanir.
	var groupID int
	//Esleme yapiliyor ise eslenen Tilelarin groupID'si
	if len(grpID) > 0 {
		groupID = grpID[0]
	} else {
		//Eslenmiyor direkt aciliyor ise Global siradaki GroupID atanir.
		groupID = Game.GenerateGroupID()
	}
	for _, tile := range setTiles {
		tile.IsOpend = true
		tile.GroupID = &groupID
	}
}

// Taslarin icinde Acilan var mi ?
func HasOpenTail(tiles ...*Model.Tile) bool {
	for _, tile := range tiles {
		if tile.IsOpend {
			return true
		}
	}
	return false
}

// Açılan taşlar arasında 5 çift var mı?
func HasAtLeastFivePairs(opened [][]*Model.Tile) bool {
	pairCount := 0

	for _, group := range opened {
		if len(group) == 2 {
			tile1, tile2 := group[0], group[1]

			//Iclerinde Zaten Acilmis var ise Hata Doner
			if HasOpenTail(tile1, tile2) {
				return false
			}

			// Okey olan taşların değerini bulmak için CalculateTileScore kullanıyoruz
			score1 := CalculateTileScore(tile1, 0, group, false)
			score2 := CalculateTileScore(tile2, 1, group, false)

			// Aynı sayı mı?
			if score1 == score2 {
				// Aynı renk mi ya da okey (joker) taşı var mı?
				// Eğer ikisi de okey ise kabul edilir.
				// Ya da renkleri aynı ise kabul edilir.
				if tile1.IsOkey || tile2.IsOkey || tile1.Color == tile2.Color {
					pairCount++
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	var result = pairCount >= 5
	if result {
		//Acilan Taslar Opened olarak isaretlenir.
		SetOpentiles(opened)
	}
	return result
}

// Açılan 5 cift taşa yeni pair eklenirken valid mi diye bakmak. Zaten isOpend : true ve GroupID leri var!
func HasAtLeastFivePairsForSetNewPair(opened [][]*Model.Tile) bool {
	pairCount := 0

	for _, group := range opened {
		if len(group) == 2 {
			tile1, tile2 := group[0], group[1]

			// Okey olan taşların değerini bulmak için CalculateTileScore kullanıyoruz
			score1 := CalculateTileScore(tile1, 0, group, false)
			score2 := CalculateTileScore(tile2, 1, group, false)

			// Aynı sayı mı?
			if score1 == score2 {
				// Aynı renk mi ya da okey (joker) taşı var mı?
				// Eğer ikisi de okey ise kabul edilir.
				// Ya da renkleri aynı ise kabul edilir.
				if tile1.IsOkey || tile2.IsOkey || tile1.Color == tile2.Color {
					pairCount++
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	var result = pairCount >= 5
	return result
}

// Rakipe Tas isleme. En fazla 2 tas ekliye bilirsin.
func CanAddTilesToSet(set []*Model.Tile, tiles ...*Model.Tile) bool {
	// Eğer eklemek istenen taş yoksa,eklenecek set yok ise, 2 den fazla tas eklenmek istendiginde ve cifte tas eklenmeye calisildiginda
	if len(tiles) == 0 || len(set) == 0 || len(tiles) > 2 || len(set) < 3 {
		return false
	}

	//Eklenecek Taslar Acilmis ise Hata Verir
	if HasOpenTail(tiles...) {
		return false
	}

	// Yeni taşları mevcut sete ekle
	newSet := append([]*Model.Tile{}, set...) // set'in kopyası
	newSet = append(newSet, tiles...)         // taşları ekle

	// Yeni set geçerli bir Group veya Sequence oluyor mu?
	var result = IsValidGroupOrRun(newSet)
	if result {
		//Islenen Taslar Acilmis olunur..
		SetOpenPairtiles(newSet, *set[0].GroupID)
	}
	return result
}

//5 Cift acmis kullaniciya, nasil isleme cift acildigini kontrol edecegiz ?
//Kullanicinin elindeki taslara bakip hepsi cift ise biz de cift ekleniyor mu diye bakabiliriz ?

// Cifte islenecek taslar uygun mu ?
func IsValidPair(tiles []*Model.Tile) bool {
	if len(tiles) != 2 {
		return false
	}

	//Eklenecek Taslar Acilmis ise Hata Verir
	if HasOpenTail(tiles...) {
		return false
	}

	tile1, tile2 := tiles[0], tiles[1]

	score1 := CalculateTileScore(tile1, 0, tiles, false)
	score2 := CalculateTileScore(tile2, 1, tiles, false)

	if score1 != score2 {
		return false
	}

	// Aynı sayı varsa, ya renkler aynı olacak ya da en az biri okey olmalı
	return tile1.IsOkey || tile2.IsOkey || tile1.Color == tile2.Color
}

// En az 5 Cift acmis kullaniciya cift tas isleme
func CanAddPairToPairSets(remaining []*Model.Tile, pairSets [][]*Model.Tile) bool {
	if IsValidPair(remaining) {
		//var result = HasAtLeastFivePairs(pairSets)
		var result = HasAtLeastFivePairsForSetNewPair(pairSets)
		if result {
			//Islene Pair IsOpened = true olarak isaretlenir..Ayrica yeni GroupID atanir.
			SetOpenPairtiles(remaining)
		}
		return result
	}
	return false
}

//**************************************

// Atilan bir tas rakibin herhangi bir setine islenebiliyor mu ?
func CanThrowingTileBeAddedToOpponentSets(newPair *Model.Tile, opponentSets [][]*Model.Tile) bool {
	if newPair.IsOpend {
		//Acilan tas bir daha atilamaz.
		return false
	}
	for _, set := range opponentSets {
		if CanAddTilesToSet(set, newPair) {
			return true
		}
	}
	return false
}

func SplitTilesByValidGroupsOrRuns_Old(tiles []*Model.Tile) ([][]*Model.Tile, []*Model.Tile) {
	sort.Slice(tiles, func(i, j int) bool {
		if tiles[i].Color == tiles[j].Color {
			return tiles[i].Number < tiles[j].Number
		}
		return tiles[i].Color < tiles[j].Color
	})

	n := len(tiles)
	type candidate struct {
		Indices   []int
		Group     []*Model.Tile
		OkeyCount int
	}
	var candidates []candidate

	// Tüm 3 ve daha uzun kombinasyonları oluştur
	for size := 3; size <= n; size++ {
		indexes := make([]int, size)
		var generate func(int, int)
		generate = func(start, depth int) {
			if depth == size {
				var group []*Model.Tile
				var okeyCount int
				for _, idx := range indexes {
					tile := tiles[idx]
					group = append(group, tile)
					if tile.IsOkey {
						okeyCount++
					}
				}
				if okeyCount > 2 {
					return
				}
				nonOkeys := filterNonOkeys(group)
				if isGroup(nonOkeys, okeyCount) || isSequence(nonOkeys, okeyCount) {
					tmp := make([]*Model.Tile, len(group))
					copy(tmp, group)
					tmpIdx := make([]int, len(indexes))
					copy(tmpIdx, indexes)
					candidates = append(candidates, candidate{
						Indices:   tmpIdx,
						Group:     tmp,
						OkeyCount: okeyCount,
					})
				}
				return
			}
			for i := start; i <= n-(size-depth); i++ {
				indexes[depth] = i
				generate(i+1, depth+1)
			}
		}
		generate(0, 0)
	}

	// Okey sayısına göre azdan çoğa sırala
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].OkeyCount < candidates[j].OkeyCount
	})

	used := make([]bool, n)
	var result [][]*Model.Tile

	for _, cand := range candidates {
		conflict := false
		for _, idx := range cand.Indices {
			if used[idx] {
				conflict = true
				break
			}
		}
		if !conflict {
			for _, idx := range cand.Indices {
				used[idx] = true
			}
			result = append(result, cand.Group)
		}
	}

	var remaining []*Model.Tile
	for i, u := range used {
		if !u {
			remaining = append(remaining, tiles[i])
		}
	}

	return result, remaining
}

func SplitTilesByValidGroupsOrRuns_X(tiles []*Model.Tile) ([][]*Model.Tile, []*Model.Tile) {
	sort.Slice(tiles, func(i, j int) bool {
		if tiles[i].Color == tiles[j].Color {
			return tiles[i].Number > tiles[j].Number // Büyükten küçüğe sayı sıralaması
		}
		return tiles[i].Color < tiles[j].Color // Renk sıralaması aynı kalıyor (küçükten büyüğe)
	})

	n := len(tiles)
	type candidate struct {
		Indices   []int
		Group     []*Model.Tile
		OkeyCount int
	}
	var candidates []candidate

	//for size := 3; size <= n; size++ {
	for size := n; size >= 3; size-- {
		indexes := make([]int, size)
		var generate func(int, int)
		generate = func(start, depth int) {
			if depth == size {
				var group []*Model.Tile
				var okeyCount int
				for _, idx := range indexes {
					tile := tiles[idx]
					group = append(group, tile)
					if tile.IsOkey {
						okeyCount++
					}
				}
				if okeyCount > 2 {
					return
				}
				nonOkeys := filterNonOkeys(group)
				if isGroup(nonOkeys, okeyCount) || isSequence(nonOkeys, okeyCount) {
					tmp := make([]*Model.Tile, len(group))
					copy(tmp, group)
					tmpIdx := make([]int, len(indexes))
					copy(tmpIdx, indexes)
					candidates = append(candidates, candidate{
						Indices:   tmpIdx,
						Group:     tmp,
						OkeyCount: okeyCount,
					})
				}
				return
			}
			for i := start; i <= n-(size-depth); i++ {
				indexes[depth] = i
				generate(i+1, depth+1)
			}
		}
		generate(0, 0)
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].OkeyCount < candidates[j].OkeyCount
	})

	used := make([]bool, n)
	var result [][]*Model.Tile

	for _, cand := range candidates {
		conflict := false
		for _, idx := range cand.Indices {
			if used[idx] {
				conflict = true
				break
			}
		}
		if !conflict {
			for _, idx := range cand.Indices {
				used[idx] = true
			}
			result = append(result, cand.Group)
		}
	}

	for _, g := range result {
		sort.Slice(g, func(i, j int) bool {
			numI := getEffectiveNumber(g[i], g)
			numJ := getEffectiveNumber(g[j], g)

			if g[i].Color == g[j].Color {
				return numI < numJ
			}
			return g[i].Color < g[j].Color
		})
	}

	// Kullanılmayan taşlar kalanlara eklenir
	var remaining []*Model.Tile
	for i, u := range used {
		if !u {
			remaining = append(remaining, tiles[i])
		}
	}

	return result, remaining
}

func SplitTilesByValidGroupsOrRuns_XX(tiles []*Model.Tile) ([][]*Model.Tile, []*Model.Tile) {
	n := len(tiles)

	// Öncelikle taşları renk ve sayı bazında sıralayalım
	sort.Slice(tiles, func(i, j int) bool {
		if tiles[i].Color == tiles[j].Color {
			return tiles[i].Number < tiles[j].Number // küçükten büyüğe
		}
		return tiles[i].Color < tiles[j].Color
	})

	type candidate struct {
		Indices   []int
		Group     []*Model.Tile
		OkeyCount int
	}

	// Tüm geçerli grupları bul
	var allGroups []candidate

	for size := 3; size <= n; size++ {
		indexes := make([]int, size)
		var generate func(int, int)
		generate = func(start, depth int) {
			if depth == size {
				var group []*Model.Tile
				okeyCount := 0
				for _, idx := range indexes {
					t := tiles[idx]
					group = append(group, t)
					if t.IsOkey {
						okeyCount++
					}
				}
				if okeyCount > 2 {
					return
				}
				nonOkeys := filterNonOkeys(group)
				if isGroup(nonOkeys, okeyCount) || isSequence(nonOkeys, okeyCount) {
					tmp := make([]*Model.Tile, len(group))
					copy(tmp, group)
					tmpIdx := make([]int, len(indexes))
					copy(tmpIdx, indexes)
					allGroups = append(allGroups, candidate{
						Indices:   tmpIdx,
						Group:     tmp,
						OkeyCount: okeyCount,
					})
				}
				return
			}
			for i := start; i <= n-(size-depth); i++ {
				indexes[depth] = i
				generate(i+1, depth+1)
			}
		}
		generate(0, 0)
	}

	// Artık elimizde tüm geçerli gruplar var.

	// Burada "grupların kombinasyonlarını" deneyip,
	// Maksimum taş açan, maksimum puanlı kombinasyonu bulacağız.

	// Bu, backtracking / bitmask ile yapılabilir.

	maxTilesUsed := 0
	maxScore := 0
	var bestCombination [][]*Model.Tile

	var backtrack func(start int, used map[int]bool, current [][]*Model.Tile)
	backtrack = func(start int, used map[int]bool, current [][]*Model.Tile) {
		// Kullanılan taş sayısı ve puanı hesapla
		usedCount := len(used)
		score := 0
		for _, group := range current {
			score += sumGroupScore(group)
		}

		// Eğer daha çok taş kullanıyorsak ya da aynı taş sayısında daha yüksek puan varsa güncelle
		if usedCount > maxTilesUsed || (usedCount == maxTilesUsed && score > maxScore) {
			maxTilesUsed = usedCount
			maxScore = score
			bestCombination = make([][]*Model.Tile, len(current))
			for i := range current {
				bestCombination[i] = make([]*Model.Tile, len(current[i]))
				copy(bestCombination[i], current[i])
			}
		}

		for i := start; i < len(allGroups); i++ {
			canUse := true
			for _, idx := range allGroups[i].Indices {
				if used[idx] {
					canUse = false
					break
				}
			}
			if !canUse {
				continue
			}

			// Kullanılan taşları işaretle
			for _, idx := range allGroups[i].Indices {
				used[idx] = true
			}
			current = append(current, allGroups[i].Group)

			backtrack(i+1, used, current)

			// Geri al
			current = current[:len(current)-1]
			for _, idx := range allGroups[i].Indices {
				delete(used, idx)
			}
		}
	}

	backtrack(0, make(map[int]bool), [][]*Model.Tile{})

	// Kullanılmayan taşlar
	used := make(map[int]bool)
	for _, group := range bestCombination {
		for _, t := range group {
			for i, tile := range tiles {
				if tile == t {
					used[i] = true
					break
				}
			}
		}
	}
	var remaining []*Model.Tile
	for i, tile := range tiles {
		if !used[i] {
			remaining = append(remaining, tile)
		}
	}

	// Her grubu sayısal ve renk olarak sıralayalım
	/*for _, g := range bestCombination {
		sort.Slice(g, func(i, j int) bool {
			numI := getEffectiveNumber(g[i], g)
			numJ := getEffectiveNumber(g[j], g)
			if g[i].Color == g[j].Color {
				return numI < numJ
			}
			return g[i].Color < g[j].Color
		})
	}*/

	//Okey siralamasinda yukaridaki loop hatali calisiyordu!
	for i, g := range bestCombination {
		bestCombination[i] = sortGroupByEffectiveNumber(g)
	}

	return bestCombination, remaining
}

func SplitTilesByValidGroupsOrRuns(tiles []*Model.Tile) ([][]*Model.Tile, []*Model.Tile) {
	n := len(tiles)

	// Renk ve Sayıya göre sırala => Renge gore grupla ve sirala sonra Sayiya gore sirala
	sort.Slice(tiles, func(i, j int) bool {
		if tiles[i].Color == tiles[j].Color {
			return tiles[i].Number < tiles[j].Number
		}
		return tiles[i].Color < tiles[j].Color
	})

	//Test All Tiles Sorting Data
	/*for _, tile := range tiles {
		fmt.Printf("Base Tile List(Number=%v, Color=%v)\n", tile.Number, GetEnumName(ColorEnum, tile.Color))
	}
	fmt.Println("-----------------------")*/

	type candidate struct {
		Indices []int
		Group   []*Model.Tile
	}

	var allGroups []candidate

	for size := 3; size <= n; size++ {
		indices := make([]int, size)
		var generate func(start, depth int)
		generate = func(start, depth int) {
			if depth == size { //Bu gerceklendiginde bir kombinasyon tamamlanmış oluyor.
				group := make([]*Model.Tile, size)
				okeyCount := 0
				for i, idx := range indices {
					tile := tiles[idx]
					group[i] = tile
					if tile.IsOkey {
						okeyCount++
					}
				}

				//Test Check All Combinations for 3 Depth
				/*if depth == 3 {
					for _, tile := range group {
						fmt.Printf("tile group(Number=%v, Color=%v)\n", tile.Number, GetEnumName(ColorEnum, tile.Color))
					}
					fmt.Println("-----------------------")
				}*/
				//------------------------

				if okeyCount > 2 { //Eğer grupta 2'den fazla okey varsa, bu kombinasyonu geçersiz sayıyoruz.
					return
				}
				nonOkeys := filterNonOkeys(group)
				if isGroup(nonOkeys, okeyCount) || isSequence(nonOkeys, okeyCount) {
					tmp := make([]*Model.Tile, len(group))
					copy(tmp, group)
					tmpIdx := make([]int, len(indices))
					copy(tmpIdx, indices)
					allGroups = append(allGroups, candidate{
						Indices: tmpIdx,
						Group:   tmp,
					})
				}
				return
			}
			//Tum kombinasyonlarin alinmasi icin Recursive olarak bir sonraki derinlik ile cagrilir
			for i := start; i <= n-(size-depth); i++ {
				indices[depth] = i
				generate(i+1, depth+1)
			}
		}
		generate(0, 0)
	}

	// Backtracking ile en iyi kombinasyonu bul
	var (
		maxTilesUsed    int
		maxScore        int
		bestCombination [][]*Model.Tile
	)

	// En fazla taş kullanılan ve en yüksek skoru veren kombinasyon aranıyor.

	var backtrack func(start int, used map[int]bool, current [][]*Model.Tile)
	backtrack = func(start int, used map[int]bool, current [][]*Model.Tile) {

		// Şu ana kadarki score hesaplaniyor.
		usedCount := len(used)
		score := 0
		for _, group := range current {
			score += sumGroupScore(group)
		}

		//Eğer bu çözüm daha iyi ise bestCombination olarak kaydediliyor.
		if usedCount > maxTilesUsed || (usedCount == maxTilesUsed && score > maxScore) {
			maxTilesUsed = usedCount
			maxScore = score
			bestCombination = deepCopyGroups(current)
		}

		//Her grup deneniyor..
		for i := start; i < len(allGroups); i++ {
			canUse := true
			for _, idx := range allGroups[i].Indices {
				//Aynı taş birden fazla kombinasyonda kullanılamaz.
				if used[idx] {
					canUse = false
					break
				}
			}
			if !canUse {
				continue
			}

			//Taşları işaretle ve geri çağır
			for _, idx := range allGroups[i].Indices {
				used[idx] = true
			}
			backtrack(i+1, used, append(current, allGroups[i].Group))
			for _, idx := range allGroups[i].Indices {
				delete(used, idx)
			}
		}
	}

	backtrack(0, make(map[int]bool), [][]*Model.Tile{})

	// Kullanılmayan taşları belirle
	usedIndices := make(map[*Model.Tile]bool)
	for _, group := range bestCombination {
		for _, tile := range group {
			usedIndices[tile] = true
		}
	}

	var remaining []*Model.Tile
	for _, tile := range tiles {
		if !usedIndices[tile] {
			remaining = append(remaining, tile)
		}
	}

	// Grupları etkili sayıya göre sırala
	//Her grup içindeki taşlar, gerçek temsil ettikleri sayıya göre sıralanır (özellikle Okey için yazdim).
	for i, g := range bestCombination {
		bestCombination[i] = sortGroupByEffectiveNumber(g)
	}

	return bestCombination, remaining
}

func deepCopyGroups(groups [][]*Model.Tile) [][]*Model.Tile {
	copied := make([][]*Model.Tile, len(groups))
	for i := range groups {
		copied[i] = make([]*Model.Tile, len(groups[i]))
		copy(copied[i], groups[i])
	}
	return copied
}

// sumGroupScore grubun toplam puanını hesaplar (toplam taş sayıları ya da puanlarını döner)
// Okey/joker etkisi varsa CalculateTileScore ile yönetilir.
func sumGroupScore(group []*Model.Tile) int {
	isSeq := isSequence(filterNonOkeys(group), countOkeys(group))
	total := 0
	for i, tile := range group {
		total += CalculateTileScore(tile, i, group, isSeq)
	}
	return total
}

func sumAllGroupsNumbers(groups [][]*Model.Tile) int {
	total := 0
	for _, group := range groups {
		isSeq := isSequence(filterNonOkeys(group), countOkeys(group))
		for index, tile := range group {

			total += CalculateTileScore(tile, index, group, isSeq)
		}
	}
	return total
}

// Taş Okey ise, hangi numarayı temsil ettiğini gruba göre çözümler.
func getEffectiveNumber(tile *Model.Tile, group []*Model.Tile) int {
	//if tile.IsOkey || tile.IsJoker {
	if tile.IsOkey {
		// Joker veya Okey'in neyi temsil ettiğini gruba bakarak çözümle
		// Basit yaklaşım: Grup içindeki en yaygın (veya eksik) değeri bul
		nonOkeys := filterNonOkeys(group)

		if isGroup(nonOkeys, countOkeys(group)) {
			// Grup ise: Aynı sayı, farklı renkler
			for _, t := range nonOkeys {
				return t.Number
			}
		} else if isSequence(nonOkeys, countOkeys(group)) {
			// Seri ise: sırayla artan sayılar aynı renk
			// Eksik olan sayıyı bul
			nums := []int{}
			for _, t := range nonOkeys {
				nums = append(nums, t.Number)
			}
			sort.Ints(nums)
			expected := nums[0]
			for _, n := range nums {
				if n != expected {
					return expected // eksik olan burası
				}
				expected++
			}
			return expected // Son eksik olan sayı
		}
	}
	// Normal taş ise kendi sayısı
	return tile.Number
}

// Okey duzgun siralanmayinca yazdim. Her taşı getEffectiveNumber ile değerleyip sıraya dizer.
func sortGroupByEffectiveNumber(group []*Model.Tile) []*Model.Tile {
	type tileWithValue struct {
		tile  *Model.Tile
		value int
	}

	var withValues []tileWithValue
	for _, t := range group {
		withValues = append(withValues, tileWithValue{
			tile:  t,
			value: getEffectiveNumber(t, group),
		})
	}

	sort.SliceStable(withValues, func(i, j int) bool {
		return withValues[i].value < withValues[j].value
	})

	var sorted []*Model.Tile
	for _, tw := range withValues {
		sorted = append(sorted, tw.tile)
	}
	return sorted
}

func SplitTilesByValidPairs_Old(tiles []*Model.Tile) ([][]*Model.Tile, []*Model.Tile) {
	bestCombination := make([][]*Model.Tile, 0)
	used := make(map[int]bool)

	for i := 0; i < len(tiles); i++ {
		if used[i] {
			continue
		}

		for j := i + 1; j < len(tiles); j++ {
			if used[j] {
				continue
			}

			t1 := tiles[i]
			t2 := tiles[j]

			// Okey her taşla eşleşebilir
			if t1.IsOkey || t2.IsOkey {
				bestCombination = append(bestCombination, []*Model.Tile{t1, t2})
				used[i], used[j] = true, true
				break
			}

			// Joker veya normal taş: renk ve sayı eşleşmeli
			if t1.Color == t2.Color && t1.Number == t2.Number {
				bestCombination = append(bestCombination, []*Model.Tile{t1, t2})
				used[i], used[j] = true, true
				break
			}
		}
	}

	// Geri kalan taşları toplama
	var remaining []*Model.Tile
	for i, t := range tiles {
		if !used[i] {
			remaining = append(remaining, t)
		}
	}

	return bestCombination, remaining
}

func SplitTilesByValidPairs(tiles []*Model.Tile) ([][]*Model.Tile, []*Model.Tile) {
	type pairKey struct {
		Color  int
		Number int
	}

	grouped := make(map[pairKey][]*Model.Tile)
	var okeys []*Model.Tile

	// Taşları sınıflandır
	for _, tile := range tiles {
		if tile.IsOkey {
			okeys = append(okeys, tile)
		} else {
			key := pairKey{Color: tile.Color, Number: tile.Number}
			grouped[key] = append(grouped[key], tile)
		}
	}

	var pairs [][]*Model.Tile
	used := make(map[*Model.Tile]bool)

	// Aynı renk ve sayıya sahip taşlardan çift oluştur
	for _, group := range grouped {
		available := []*Model.Tile{}
		for _, t := range group {
			if !used[t] {
				available = append(available, t)
			}
		}
		for len(available) >= 2 {
			pair := available[:2]
			pairs = append(pairs, pair)
			used[pair[0]] = true
			used[pair[1]] = true
			available = available[2:]
		}
	}

	// Kalan tek taşları okey ile eşleştir
	for _, group := range grouped {
		//Okey yok ise donguye gerek yok
		if len(okeys) == 0 {
			break
		}
		available := []*Model.Tile{}
		for _, t := range group {
			if !used[t] {
				available = append(available, t)
			}
		}
		for len(available) >= 1 && len(okeys) > 0 {
			t := available[0]
			okey := okeys[0]
			pairs = append(pairs, []*Model.Tile{t, okey})
			used[t] = true
			used[okey] = true
			available = available[1:]
			okeys = okeys[1:]
		}
	}

	// Kalan taşlar
	var remaining []*Model.Tile
	for _, t := range tiles {
		if !used[t] {
			remaining = append(remaining, t)
		}
	}

	return pairs, remaining
}
