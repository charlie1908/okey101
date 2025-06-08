package Core

import (
	"okey101/Model"
	"sort"
)

// Kendisine gelen tas dizilimlerinin Ayni Sayi Farkjli Renk (Group) veya Sirali Ayni Renk (RUN) olma durumuna bakma.
func IsValidGroupOrRun(tiles []Model.Tile) bool {
	if len(tiles) < 3 {
		return false
	}

	okeyCount := countOkeys(tiles)
	nonOkeyTiles := filterNonOkeys(tiles)

	return isGroup(nonOkeyTiles, okeyCount) || isSequence(nonOkeyTiles, okeyCount)
}

// Kirmizi 5, Mavi 5, Sari 5
func isGroup(tiles []Model.Tile, okeyCount int) bool {
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
func isSequence(tiles []Model.Tile, okeyCount int) bool {
	if len(tiles) == 0 {
		return false
	}

	//Ayni renkte degiller ise zaten seri olamaz
	if !allSameColor(tiles) {
		return false
	}

	//Kullanici o seriyi sirali vermese bile siralar dogru mu diye siralayip bakariz..
	sort.Slice(tiles, func(i, j int) bool {
		return tiles[i].Number < tiles[j].Number
	})

	//Arada bosluklar var ise gerekli Okey sayisina bakilir..
	neededOkeys := calculateNeededOkeysForRun(tiles, okeyCount)

	//Ihtiyac duyulan OkeyCount eldekinden fazla ise veya toplam serideki tas sayisi 3'den kucuk ise false firlatilir.
	return neededOkeys <= okeyCount && len(tiles)+okeyCount >= 3
}

func calculateNeededOkeysForRun(tiles []Model.Tile, maxOkeys int) int {
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

func countOkeys(tiles []Model.Tile) int {
	count := 0
	for _, tile := range tiles {
		if tile.IsOkey {
			count++
		}
	}
	return count
}

func filterNonOkeys(tiles []Model.Tile) []Model.Tile {
	result := make([]Model.Tile, 0, len(tiles))
	for _, tile := range tiles {
		if !tile.IsOkey {
			result = append(result, tile)
		}
	}
	return result
}

func allSameColor(tiles []Model.Tile) bool {
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

func CalculateTileScore(tile Model.Tile, index int, tiles []Model.Tile, isSequence bool) int {
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
			nonOkeys = append(nonOkeys, t)
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
func CanOpenTiles(opened [][]Model.Tile) bool {
	totalScore := 0
	for _, group := range opened {

		//Once gecerliligi kontrol et!
		if !IsValidGroupOrRun(group) {
			return false
		}

		isSeq := isSequence(filterNonOkeys(group), countOkeys(group))
		for index, tile := range group {

			totalScore += CalculateTileScore(tile, index, group, isSeq)
		}
	}
	return totalScore >= 101
}

// Açılan taşlar arasında 5 çift var mı?
func HasAtLeastFivePairs(opened [][]Model.Tile) bool {
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

	return pairCount >= 5
}

// Rakipe Tas isleme. En fazla 2 tas ekliye bilirsin.
func CanAddTilesToSet(set []Model.Tile, tiles ...Model.Tile) bool {
	// Eğer eklemek istenen taş yoksa,eklenecek set yok ise, 2 den fazla tas eklenmek istendiginde ve cifte tas eklenmeye calisildiginda
	if len(tiles) == 0 || len(set) == 0 || len(tiles) > 2 || len(set) < 3 {
		return false
	}

	// Yeni taşları mevcut sete ekle
	newSet := append([]Model.Tile{}, set...) // set'in kopyası
	newSet = append(newSet, tiles...)        // taşları ekle

	// Yeni set geçerli bir Group veya Sequence oluyor mu?
	return IsValidGroupOrRun(newSet)
}

//5 Cift acmis kullaniciya, nasil isleme cift acildigini kontrol edecegiz ?
//Kullanicinin elindeki taslara bakip hepsi cift ise biz de cift ekleniyor mu diye bakabiliriz ?

// Cifte islenecek taslar uygun mu ?
func IsValidPair(tiles []Model.Tile) bool {
	if len(tiles) != 2 {
		return false
	}

	tile1, tile2 := tiles[0], tiles[1]

	score1 := CalculateTileScore(tile1, 0, tiles, false)
	score2 := CalculateTileScore(tile2, 1, tiles, false)

	if score1 != score2 {
		return false
	}

	// Aynı sayı varsa, ya renkler aynı olacak ya da en az biri okey/joker olmalı
	return tile1.IsOkey || tile2.IsOkey || tile1.Color == tile2.Color
}

// En az 5 Cift acmis kullaniciya cift tas isleme
func CanAddPairToPairSets(remaining []Model.Tile, pairSets [][]Model.Tile) bool {
	if IsValidPair(remaining) {
		return HasAtLeastFivePairs(pairSets)
	}
	return false
}

//**************************************

// Atilan bir tas rakibin herhangi bir setine islenebiliyor mu ?
func CanThrowingTileBeAddedToOpponentSets(newPair Model.Tile, opponentSets [][]Model.Tile) bool {
	for _, set := range opponentSets {
		if CanAddTilesToSet(set, newPair) {
			return true
		}
	}
	return false
}
