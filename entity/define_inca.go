package entity

import (
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

const (
	ColsIncaMatrix = 5
	RowsIncaMatrix = 3
)

var IncaLowSymbols = []pb.SiXiangSymbol{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_J,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_Q,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_K,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_A,
}

var IncaMidSymbols = []pb.SiXiangSymbol{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_SUN,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_EAGLE_GARUDA,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_ANTIQUE,
}

var IncaHighSymbols = []pb.SiXiangSymbol{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_HEARTS,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_DIAMONDS,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_SPADES,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_CLUBS,
}
var (
	IncalAllSymbol      []pb.SiXiangSymbol
	IncaPaylineMap      = orderedmap.New[int, []int]()
	IncaRatioPaylineMap = make(map[pb.SiXiangSymbol]map[int32]float64)
)

func init() {
	{
		IncalAllSymbol = append(IncalAllSymbol, IncaLowSymbols...)
		IncalAllSymbol = append(IncalAllSymbol, IncaMidSymbols...)
		IncalAllSymbol = append(IncalAllSymbol, IncaHighSymbols...)
		IncalAllSymbol = append(IncalAllSymbol,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER,
		)
	}
	{
		var m = map[int32]float64{2: 0, 3: 5, 4: 25, 5: 125}
		IncaRatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_J] = m
		IncaRatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_Q] = m
		IncaRatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_K] = m
	}
	{
		var m = map[int32]float64{2: 0, 3: 10, 4: 50, 5: 250}
		IncaRatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_SPADES] = m
		IncaRatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_HEARTS] = m
		IncaRatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_A] = m
	}
	{
		var m = map[int32]float64{2: 0, 3: 15, 4: 75, 5: 500}
		IncaRatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_DIAMONDS] = m
		IncaRatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_CLUBS] = m
	}
	{
		var m = map[int32]float64{2: 2, 3: 25, 4: 125, 5: 725}
		IncaRatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_EAGLE_GARUDA] = m
		IncaRatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_ANTIQUE] = m
	}
	{
		var m = map[int32]float64{2: 2, 3: 50, 4: 250, 5: 1250}
		m[pb.SiXiangSymbol_SI_XIANG_SYMBOL_SUN] = m
	}
	{
		var m = map[int32]float64{5: 1000}
		IncaRatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD] = m
	}
	{
		var m = map[int32]float64{2: 40, 3: 50}
		IncaRatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER] = m
	}
	{
		idx := 1
		IncaPaylineMap.Set(idx, []int{5, 6, 7, 8, 9})
		idx = 2
		IncaPaylineMap.Set(idx, []int{0, 1, 2, 3, 4})
		idx = 3
		IncaPaylineMap.Set(idx, []int{10, 11, 12, 13, 14})
		idx = 4
		IncaPaylineMap.Set(idx, []int{0, 6, 12, 8, 4})
		idx = 5
		IncaPaylineMap.Set(idx, []int{10, 6, 2, 8, 14})
		idx = 6
		IncaPaylineMap.Set(idx, []int{0, 1, 7, 13, 14})
		idx = 7
		IncaPaylineMap.Set(idx, []int{10, 11, 7, 3, 4})
		idx = 8
		IncaPaylineMap.Set(idx, []int{5, 1, 7, 13, 9})
		idx = 9
		IncaPaylineMap.Set(idx, []int{5, 11, 7, 3, 9})
		idx = 10
		IncaPaylineMap.Set(idx, []int{5, 1, 2, 8, 4})
		idx = 11
		IncaPaylineMap.Set(idx, []int{5, 11, 12, 8, 14})
		idx = 12
		IncaPaylineMap.Set(idx, []int{0, 6, 2, 3, 9})
		idx = 13
		IncaPaylineMap.Set(idx, []int{10, 6, 12, 13, 9})
		idx = 14
		IncaPaylineMap.Set(idx, []int{0, 11, 2, 13, 4})
		idx = 15
		IncaPaylineMap.Set(idx, []int{10, 1, 12, 3, 14})
		idx = 16
		IncaPaylineMap.Set(idx, []int{5, 1, 12, 3, 9})
		idx = 17
		IncaPaylineMap.Set(idx, []int{5, 11, 2, 13, 9})
		idx = 18
		IncaPaylineMap.Set(idx, []int{0, 6, 7, 8, 4})
		idx = 19
		IncaPaylineMap.Set(idx, []int{10, 6, 7, 8, 14})
		idx = 20
		IncaPaylineMap.Set(idx, []int{0, 11, 12, 13, 4})
	}
}
