package entity

import (
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

const (
	ColsJuicyMatrix  = 5
	RowsJuicynMatrix = 3
)

var JuicyLowSymbols = []pb.SiXiangSymbol{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_J,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_Q,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_K,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_A,
}

var JuicyMidSymbol = []pb.SiXiangSymbol{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STRAWBERRY,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_WATERMELON,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_PINAPPLE,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_MANGOSTEEN,
}

var JuicyHighSymbol = []pb.SiXiangSymbol{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_DIAMOND,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_GREEN,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_VIOLET,
}

var JuicySpecialSymbol = []pb.SiXiangSymbol{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_SPIN,
}

var JuicyFruitBasektSymbol = []pb.SiXiangSymbol{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_GRAND,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MAJOR,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MINOR,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MINI,
}

var JuicyFruitRainSybol = []pb.SiXiangSymbol{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_SPIN,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_1,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_2,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_3,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_4,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_5,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_6,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_7,
}

func IsFruitBasketSymbol(symbol pb.SiXiangSymbol) bool {
	_, exist := JuicyBasketSymbol[symbol]
	return exist
}

func IsFruitJPSymbol(symbol pb.SiXiangSymbol) bool {
	if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_GRAND ||
		symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MAJOR ||
		symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MINOR ||
		symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MINI {
		return true
	}
	return false
}

var JuicyBasketSymbol = map[pb.SiXiangSymbol]SymbolInfo{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_SPIN: {
		NumOccur: 0,
		Value:    Range{Min: 0, Max: 0},
	},

	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_1: {
		NumOccur: 3,
		Value:    Range{Min: 10, Max: 15},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_2: {
		NumOccur: 3,
		Value:    Range{Min: 20, Max: 30},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_3: {
		NumOccur: 3,
		Value:    Range{Min: 50, Max: 70},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_4: {
		NumOccur: 2,
		Value:    Range{Min: 100, Max: 120},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_5: {
		NumOccur: 2,
		Value:    Range{Min: 180, Max: 200},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_6: {
		NumOccur: 1,
		Value:    Range{Min: 250, Max: 300},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_7: {
		NumOccur: 1,
		Value:    Range{Min: 250, Max: 300},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MINI: {
		NumOccur: 0,
		Value:    Range{Min: 50, Max: 50},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MINOR: {
		NumOccur: 0,
		Value:    Range{Min: 100, Max: 100},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MAJOR: {
		NumOccur: 0,
		Value:    Range{Min: 500, Max: 500},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_GRAND: {
		NumOccur: 0,
		Value:    Range{Min: 5000, Max: 5000},
	},
}

var JuiceAllSymbols []pb.SiXiangSymbol
var JuiceAllSymbolsWildRatio1_2 []pb.SiXiangSymbol
var JuiceAllSymbolsWildRatio1_5 []pb.SiXiangSymbol
var JuiceAllSymbolsWildRatio2_0 []pb.SiXiangSymbol

var MapJuicyPaylineIdx = orderedmap.New[int, []int]()

var RatioJuicyPaylineMap map[pb.SiXiangSymbol]map[int32]float64

func init() {
	for i := 0; i < 5; i++ {

		JuiceAllSymbols = append(JuiceAllSymbols, JuicyLowSymbols...)
		JuiceAllSymbols = append(JuiceAllSymbols, JuicyMidSymbol...)
		JuiceAllSymbols = append(JuiceAllSymbols, JuicyHighSymbol...)
	}
	JuiceAllSymbols = append(JuiceAllSymbols, JuicySpecialSymbol...)
	list := make([]pb.SiXiangSymbol, 0, len(JuiceAllSymbols)*10)

	for i := 0; i < 10; i++ {
		list = append(list, JuiceAllSymbols...)
	}
	// JuiceAllSymbols = list

	JuiceAllSymbolsWildRatio1_2 = append(JuiceAllSymbolsWildRatio1_2, list...)
	JuiceAllSymbolsWildRatio1_2 = append(JuiceAllSymbolsWildRatio1_2, SliceRepeat(2, pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD)...)

	JuiceAllSymbolsWildRatio1_5 = append(JuiceAllSymbolsWildRatio1_5, list...)
	JuiceAllSymbolsWildRatio1_5 = append(JuiceAllSymbolsWildRatio1_5, SliceRepeat(5, pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD)...)

	JuiceAllSymbolsWildRatio2_0 = append(JuiceAllSymbolsWildRatio2_0, list...)
	JuiceAllSymbolsWildRatio2_0 = append(JuiceAllSymbolsWildRatio2_0, SliceRepeat(10, pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD)...)

	idx := 1
	MapJuicyPaylineIdx.Set(idx, []int{5, 6, 7, 8, 9})
	idx++
	MapJuicyPaylineIdx.Set(idx, []int{0, 1, 2, 3, 4})
	idx++
	MapJuicyPaylineIdx.Set(idx, []int{10, 11, 12, 13, 14})
	idx++
	MapJuicyPaylineIdx.Set(idx, []int{0, 6, 12, 8, 4})
	idx++
	MapJuicyPaylineIdx.Set(idx, []int{10, 6, 2, 8, 14})

	idx++
	MapJuicyPaylineIdx.Set(idx, []int{0, 1, 7, 13, 14})
	idx++
	MapJuicyPaylineIdx.Set(idx, []int{10, 11, 7, 3, 4})
	idx++
	MapJuicyPaylineIdx.Set(idx, []int{5, 1, 7, 13, 9})
	idx++
	MapJuicyPaylineIdx.Set(idx, []int{5, 11, 7, 3, 9})
	idx++
	MapJuicyPaylineIdx.Set(idx, []int{5, 1, 2, 8, 4})

	idx++
	MapJuicyPaylineIdx.Set(idx, []int{5, 11, 12, 8, 14})
	idx++
	MapJuicyPaylineIdx.Set(idx, []int{0, 6, 2, 3, 9})
	idx++
	MapJuicyPaylineIdx.Set(idx, []int{10, 6, 12, 13, 9})
	idx++
	MapJuicyPaylineIdx.Set(idx, []int{0, 11, 2, 13, 4})
	idx++
	MapJuicyPaylineIdx.Set(idx, []int{10, 1, 12, 3, 14})

	idx++
	MapJuicyPaylineIdx.Set(idx, []int{5, 1, 12, 3, 9})
	idx++
	MapJuicyPaylineIdx.Set(idx, []int{5, 11, 2, 13, 9})
	idx++
	MapJuicyPaylineIdx.Set(idx, []int{0, 6, 7, 8, 4})
	idx++
	MapJuicyPaylineIdx.Set(idx, []int{10, 6, 7, 8, 14})
	idx++
	MapJuicyPaylineIdx.Set(idx, []int{0, 11, 12, 13, 4})

	RatioJuicyPaylineMap = make(map[pb.SiXiangSymbol]map[int32]float64)
	{
		var m = map[int32]float64{2: 0, 3: 5, 4: 25, 5: 125}
		RatioJuicyPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_J] = m
		RatioJuicyPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_Q] = m
		RatioJuicyPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_K] = m
	}
	{
		var m = map[int32]float64{2: 0, 3: 10, 4: 55, 5: 250}
		RatioJuicyPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_A] = m
		RatioJuicyPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_MANGOSTEEN] = m
		RatioJuicyPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_PINAPPLE] = m

	}
	{
		var m = map[int32]float64{2: 0, 3: 15, 4: 75, 5: 550}
		RatioJuicyPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_WATERMELON] = m
		RatioJuicyPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STRAWBERRY] = m
	}
	{
		var m = map[int32]float64{2: 2, 3: 25, 4: 125, 5: 750}
		RatioJuicyPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_DIAMOND] = m
		RatioJuicyPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_VIOLET] = m
	}
	{
		var m = map[int32]float64{2: 2, 3: 50, 4: 250, 5: 1250}
		RatioJuicyPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_GREEN] = m
	}

}

func NewJuicyMatrix() SlotMatrix {
	m := SlotMatrix{
		List:      make([]pb.SiXiangSymbol, ColsJuicyMatrix*RowsJuicynMatrix),
		Cols:      ColsJuicyMatrix,
		Rows:      RowsJuicynMatrix,
		Size:      ColsJuicyMatrix * RowsJuicynMatrix,
		TrackFlip: make(map[int]bool),
	}
	return m
}

func NewJuicyFruitRainMaxtrix() SlotMatrix {
	m := NewSlotMatrix(3, 5)
	for sym, val := range JuicyBasketSymbol {
		m.List = append(m.List, SliceRepeat(val.NumOccur, sym)...)
	}
	m.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_7 {
			arr := []pb.SiXiangSymbol{
				pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_7,
				pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MINI,
				pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MINOR,
				pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MAJOR,
			}
			m.List[idx] = ShuffleSlice(arr)[RandomInt(0, len(arr))]
		}
	})
	return m
}

func JuicySpinSymbol(randFn func(min, max int) int, list []pb.SiXiangSymbol) pb.SiXiangSymbol {
	lenList := len(list)
	if lenList == 0 {
		return pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
	}
	idx := randFn(0, lenList)
	randSymbol := list[idx]
	return randSymbol
}

func JuicySpinSymbolToJp(sym pb.SiXiangSymbol) pb.WinJackpot {
	switch sym {
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MINI:
		return pb.WinJackpot_WIN_JACKPOT_MINI
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MINOR:
		return pb.WinJackpot_WIN_JACKPOT_MINOR
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MAJOR:
		return pb.WinJackpot_WIN_JACKPOT_MAJOR
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_GRAND:
		return pb.WinJackpot_WIN_JACKPOT_GRAND
	default:
		return pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED
	}
}

var JuiceJpRatio = map[pb.WinJackpot]int{
	pb.WinJackpot_WIN_JACKPOT_MINI:  50,
	pb.WinJackpot_WIN_JACKPOT_MINOR: 100,
	pb.WinJackpot_WIN_JACKPOT_MAJOR: 500,
	pb.WinJackpot_WIN_JACKPOT_GRAND: 5000,
}
