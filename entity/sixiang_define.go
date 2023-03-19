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
		Value:    Range{5, 7},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_2: {
		NumOccur: 1,
		Value:    Range{10, 14},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_3: {
		NumOccur: 1,
		Value:    Range{15, 18},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MINOR: {
		NumOccur: 3,
		Value:    Range{10, 10},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MAJOR: {
		NumOccur: 3,
		Value:    Range{50, 50},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MEGA: {
		NumOccur: 3,
		Value:    Range{100, 100},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GRAND: {
		NumOccur: 3,
		Value:    Range{500, 500},
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
		Value:    Range{0.1, 0.2},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM2: {
		NumOccur: 3,
		Value:    Range{0.3, 0.7},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM3: {
		NumOccur: 2,
		Value:    Range{1, 2},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM4: {
		NumOccur: 2,
		Value:    Range{3, 6},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM5: {
		NumOccur: 1,
		Value:    Range{8, 10},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MINOR: {
		NumOccur: 0,
		Value:    Range{10, 10},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MAJOR: {
		NumOccur: 0,
		Value:    Range{50, 50},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MEGA: {
		NumOccur: 0,
		Value:    Range{100, 100},
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
		Value:    Range{0.1, 0.2},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_GOLD2: {
		NumOccur: 4,
		Value:    Range{0.3, 0.7},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_GOLD3: {
		NumOccur: 4,
		Value:    Range{1, 2},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_GOLD4: {
		NumOccur: 2,
		Value:    Range{3, 6},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_GOLD5: {
		NumOccur: 1,
		Value:    Range{8, 10},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_JP_MINOR: {
		NumOccur: 0,
		Value:    Range{10, 10},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_JP_MAJOR: {
		NumOccur: 0,
		Value:    Range{50, 50},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_JP_MEGA: {
		NumOccur: 0,
		Value:    Range{100, 100},
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
		Value:    Range{3, 3},
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

func init() {
	// MapPaylineIdx = make(map[int][]int, 0)
	// 1 - 5
	idx := 0
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
}
