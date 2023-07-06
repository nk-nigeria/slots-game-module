package entity

import (
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

const (
	ColsMatrix = 5
	RowsMatrix = 3

	ColsMatrixGoldPick = 5
	RowsMatrixGoldPick = 4

	ColsMatrixRapidPay = 5
	RowsMatrixRapidPay = 5
)

var ListSymbol = []pb.SiXiangSymbol{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_BLUE_DRAGON,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_WHITE_TIGER,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_WARRIOR,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_VERMILION_BIRD,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_10,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_J,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_Q,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_K,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_A,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD,
}

var ListSymbolLuckyDraw = map[pb.SiXiangSymbol]SymbolInfo{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_1: {
		NumOccur: 1,
		Value:    Range{1.5, 2},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_2: {
		NumOccur: 1,
		Value:    Range{3, 4},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_3: {
		NumOccur: 1,
		Value:    Range{4, 5},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MINOR: {
		NumOccur: 3,
		Value:    Range{float32(pb.WinJackpot_WIN_JACKPOT_MINOR), float32(pb.WinJackpot_WIN_JACKPOT_MINOR)},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MAJOR: {
		NumOccur: 3,
		Value:    Range{float32(pb.WinJackpot_WIN_JACKPOT_MAJOR), float32(pb.WinJackpot_WIN_JACKPOT_MAJOR)},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MEGA: {
		NumOccur: 3,
		Value:    Range{float32(pb.WinJackpot_WIN_JACKPOT_MEGA), float32(pb.WinJackpot_WIN_JACKPOT_MEGA)},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GRAND: {
		NumOccur: 3,
		Value:    Range{float32(pb.WinJackpot_WIN_JACKPOT_GRAND), float32(pb.WinJackpot_WIN_JACKPOT_GRAND)},
	},
}

var ListSymbolBonusGame = map[pb.SiXiangSymbol]SymbolInfo{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_DRAGONBALL: {
		NumOccur: 1,
		Value:    Range{0, 0},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_LUCKYDRAW: {
		NumOccur: 1,
		Value:    Range{0, 0},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDPICK: {
		NumOccur: 1,
		Value:    Range{0, 0},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_RAPIDPAY: {
		NumOccur: 1,
		Value:    Range{0, 0},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDX10: {
		NumOccur: 1,
		Value:    Range{10, 10},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDX20: {
		NumOccur: 1,
		Value:    Range{20, 20},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDX30: {
		NumOccur: 1,
		Value:    Range{30, 30},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDX50: {
		NumOccur: 1,
		Value:    Range{50, 50},
	},
}

var ListSymbolDragonPearl = map[pb.SiXiangSymbol]SymbolInfo{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_LUCKMONEY: {
		NumOccur: 3,
		Value:    Range{0, 0},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM1: {
		NumOccur: 4,
		Value:    Range{0.03, 0.06},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM2: {
		NumOccur: 3,
		Value:    Range{0.1, 0.25},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM3: {
		NumOccur: 2,
		Value:    Range{0.3, 0.6},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM4: {
		NumOccur: 2,
		Value:    Range{1, 2},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM5: {
		NumOccur: 1,
		Value:    Range{2.5, 3},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MINOR: {
		NumOccur: 0,
		Value:    Range{float32(pb.WinJackpot_WIN_JACKPOT_MINOR), float32(pb.WinJackpot_WIN_JACKPOT_MINOR)},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MAJOR: {
		NumOccur: 0,
		Value:    Range{float32(pb.WinJackpot_WIN_JACKPOT_MAJOR), float32(pb.WinJackpot_WIN_JACKPOT_MAJOR)},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MEGA: {
		NumOccur: 0,
		Value:    Range{float32(pb.WinJackpot_WIN_JACKPOT_MEGA), float32(pb.WinJackpot_WIN_JACKPOT_MEGA)},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_GRAND: {
		NumOccur: 0,
		Value:    Range{float32(pb.WinJackpot_WIN_JACKPOT_GRAND), float32(pb.WinJackpot_WIN_JACKPOT_GRAND)},
	},
}

var ListEyeSiXiang = map[pb.SiXiangSymbol]SymbolInfo{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_BIRD: {
		NumOccur: 1,
		Value:    Range{1, 1},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_TIGER: {
		NumOccur: 1,
		Value:    Range{1, 1},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_WARRIOR: {
		NumOccur: 1,
		Value:    Range{2, 2},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_DRAGON: {
		NumOccur: 1,
		Value:    Range{1, 1},
	},
}

var ListSpecialGame = []pb.SiXiangGame{
	pb.SiXiangGame_SI_XIANG_GAME_BONUS,
	pb.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL,
	pb.SiXiangGame_SI_XIANG_GAME_LUCKDRAW,
	pb.SiXiangGame_SI_XIANG_GAME_RAPIDPAY,
	pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS,
}

var ListSymbolGoldPick = map[pb.SiXiangSymbol]SymbolInfo{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_TRYAGAIN: {
		NumOccur: 4,
		Value:    Range{0, 0},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_GOLD1: {
		NumOccur: 5,
		Value:    Range{0.1, 0.3},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_GOLD2: {
		NumOccur: 4,
		Value:    Range{1, 2},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_GOLD3: {
		NumOccur: 4,
		Value:    Range{3, 4},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_GOLD4: {
		NumOccur: 2,
		Value:    Range{5, 7},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_GOLD5: {
		NumOccur: 1,
		Value:    Range{8, 10},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_JP_MINOR: {
		NumOccur: 0,
		Value:    Range{float32(pb.WinJackpot_WIN_JACKPOT_MINOR), float32(pb.WinJackpot_WIN_JACKPOT_MINOR)},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_JP_MAJOR: {
		NumOccur: 0,
		Value:    Range{float32(pb.WinJackpot_WIN_JACKPOT_MAJOR), float32(pb.WinJackpot_WIN_JACKPOT_MAJOR)},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_JP_MEGA: {
		NumOccur: 0,
		Value:    Range{float32(pb.WinJackpot_WIN_JACKPOT_MEGA), float32(pb.WinJackpot_WIN_JACKPOT_MEGA)},
	},
}

var ListSymbolRapidPay = map[pb.SiXiangSymbol]SymbolInfo{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_END: {
		NumOccur: 1,
		Value:    Range{},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X2: {
		NumOccur: 1,
		Value:    Range{2, 2},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X3: {
		NumOccur: 1,
		Value:    Range{3, 3},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X4: {
		NumOccur: 1,
		Value:    Range{4, 4},
	},
}

var ListSymbolSiXiangBonusGame = []pb.SiXiangSymbol{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_SIXANGBONUS_DRAGONPEARL_GAME,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_SIXANGBONUS_LUCKYDRAW_GAME,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_SIXANGBONUS_GOLDPICK_GAME,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_SIXANGBONUS_RAPIDPAY_GAME,
}

type Range struct {
	Min float32
	Max float32
}
type SymbolInfo struct {
	NumOccur int
	Value    Range
}

var MapPaylineIdx = orderedmap.New[int, []int]()
var RellsAllowScatter = map[int]bool{0: true, 2: true, 4: true}

var RatioPaylineMap map[pb.SiXiangSymbol]map[int32]float64
var ListSymbolSpinInSixiangNormal []pb.SiXiangSymbol

func init() {
	// MapPaylineIdx = make(map[int][]int, 0)
	// 1 - 5
	idx := 1
	MapPaylineIdx.Set(idx, []int{0, 1, 2, 3, 4})
	// MapPaylineIdx[idx] = []int{0, 1, 2, 3, 4}
	idx++
	MapPaylineIdx.Set(idx, []int{5, 6, 7, 8, 9})
	idx++
	MapPaylineIdx.Set(idx, []int{10, 11, 12, 13, 14})
	idx++
	MapPaylineIdx.Set(idx, []int{10, 6, 2, 8, 14})
	idx++
	MapPaylineIdx.Set(idx, []int{0, 6, 12, 8, 4})
	idx++
	// 5 - 10
	MapPaylineIdx.Set(idx, []int{5, 11, 7, 3, 9})
	idx++
	MapPaylineIdx.Set(idx, []int{5, 11, 7, 13, 9})
	idx++
	MapPaylineIdx.Set(idx, []int{0, 6, 2, 8, 14})
	idx++
	MapPaylineIdx.Set(idx, []int{10, 6, 12, 8, 4})
	idx++
	MapPaylineIdx.Set(idx, []int{10, 6, 2, 8, 4})
	idx++

	// 11 - 15
	MapPaylineIdx.Set(idx, []int{0, 6, 12, 8, 14})
	idx++
	MapPaylineIdx.Set(idx, []int{5, 1, 7, 13, 9})
	idx++
	MapPaylineIdx.Set(idx, []int{10, 6, 12, 8, 14})
	idx++
	MapPaylineIdx.Set(idx, []int{5, 1, 7, 3, 9})
	idx++
	MapPaylineIdx.Set(idx, []int{0, 6, 2, 8, 4})
	idx++
	//16-20
	MapPaylineIdx.Set(idx, []int{5, 6, 12, 8, 9})
	idx++
	MapPaylineIdx.Set(idx, []int{0, 1, 7, 3, 4})
	idx++
	MapPaylineIdx.Set(idx, []int{10, 11, 7, 13, 14})
	idx++
	MapPaylineIdx.Set(idx, []int{5, 6, 2, 8, 9})
	idx++
	MapPaylineIdx.Set(idx, []int{0, 1, 12, 3, 4})
	idx++
	//21-25
	MapPaylineIdx.Set(idx, []int{10, 11, 2, 13, 14})
	idx++
	MapPaylineIdx.Set(idx, []int{5, 6, 7, 13, 9})
	idx++
	MapPaylineIdx.Set(idx, []int{0, 1, 2, 8, 14})
	idx++
	MapPaylineIdx.Set(idx, []int{10, 11, 12, 8, 4})
	idx++
	MapPaylineIdx.Set(idx, []int{10, 6, 7, 8, 4})

	RatioPaylineMap = make(map[pb.SiXiangSymbol]map[int32]float64)
	{
		var m = map[int32]float64{3: 0.5, 4: 2.5, 5: 5}
		RatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_10] = m
		RatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_J] = m
		RatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_Q] = m
		RatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_K] = m
	}
	{
		var m = map[int32]float64{3: 2, 4: 10, 5: 20}
		RatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_BLUE_DRAGON] = m
	}
	{
		var m = map[int32]float64{3: 1.5, 4: 7.5, 5: 15}
		RatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_WHITE_TIGER] = m
	}
	{
		var m = map[int32]float64{3: 1.2, 4: 6, 5: 12}
		RatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_WARRIOR] = m
	}
	{
		var m = map[int32]float64{3: 1, 4: 5, 5: 10}
		RatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_VERMILION_BIRD] = m
	}

	{
		ListSymbolSpinInSixiangNormal = make([]pb.SiXiangSymbol, 0)
		listSymbolExceptWilAndScatter := make([]pb.SiXiangSymbol, 0, len(ListSymbol))
		for _, sym := range ListSymbol {
			if sym == pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD || sym == pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER {
				continue
			}
			listSymbolExceptWilAndScatter = append(listSymbolExceptWilAndScatter, sym)
		}
		ListSymbolSpinInSixiangNormal = append(ListSymbolSpinInSixiangNormal, ShuffleSlice(ListSymbol)...)
		for i := 0; i < 10; i++ {
			ListSymbolSpinInSixiangNormal = append(ListSymbolSpinInSixiangNormal, ShuffleSlice(listSymbolExceptWilAndScatter)...)
		}
	}
}

func IsSixiangEyeSymbol(sym pb.SiXiangSymbol) bool {
	_, ok := ListEyeSiXiang[sym]
	return ok
}
