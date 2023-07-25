package entity

import (
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

const (
	ColsTarzanMatrix = 5
	RowsTarzanMatrix = 3
)

var TarzanSymbols = []pb.SiXiangSymbol{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_J,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_Q,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_K,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_A,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GORILLE,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_ELEPHANT,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JAGUAR,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_SNACK,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JANE,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JANE_FATHER,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_CLAYTON,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_FREE_SPIN,
	// pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_J,
	// pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_U,
	// pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_N,
	// pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_G,
	// pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_L,
	// pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_E,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN,
	// pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_FREE_SPIN,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN,
}

var TarzanLowSymbol = map[pb.SiXiangSymbol]bool{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_J:        true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_Q:        true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_K:        true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_A:        true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GORILLE:  true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_ELEPHANT: true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JAGUAR:   true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_SNACK:    true,
}

var TarzanMidSymbol = map[pb.SiXiangSymbol]bool{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GORILLE:  true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_ELEPHANT: true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JAGUAR:   true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_SNACK:    true,
}

var TarzanHighSymbol = map[pb.SiXiangSymbol]bool{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JANE:        true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JANE_FATHER: true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_CLAYTON:     true,
}

var TarzanLetterSymbol = map[pb.SiXiangSymbol]bool{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_J: true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_U: true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_N: true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_G: true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_L: true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_E: true,
}

var TarzanJungleTreasureSymbol = map[pb.SiXiangSymbol]SymbolInfo{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_MORE_TURNX2: {
		NumOccur: 2,
		Value:    Range{Min: 2, Max: 2},
	},
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_MORE_TURNX3: {
		NumOccur: 2,
		Value:    Range{Min: 3, Max: 3},
	},
	// 10 - 30 line
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_RANDOM_1: {
		NumOccur: 3,
		Value:    Range{Min: 10, Max: 30},
	},
	// 50 - 120 line
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_RANDOM_2: {
		NumOccur: 3,
		Value:    Range{Min: 50, Max: 120},
	},
	// 200 - 300 line
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_RANDOM_3: {
		NumOccur: 2,
		Value:    Range{Min: 200, Max: 300},
	},
	// 500 - 700 line
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_RANDOM_4: {
		NumOccur: 2,
		Value:    Range{Min: 500, Max: 700},
	},
	// 1000 - 1200 line
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_RANDOM_5: {
		NumOccur: 1,
		Value:    Range{Min: 1000, Max: 1200},
	},
}

var PaylineTarzanMapping = orderedmap.New[int, []int]()

func init() {
	idx := 1
	PaylineTarzanMapping.Set(idx, []int{5, 6, 7, 8, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 1, 2, 3, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 11, 12, 13, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 6, 12, 8, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 6, 2, 8, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 1, 2, 3, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 11, 12, 13, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 1, 7, 13, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 11, 7, 3, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 11, 7, 3, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 1, 7, 13, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 6, 7, 8, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 6, 7, 8, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 6, 2, 8, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 6, 12, 8, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 6, 2, 8, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 6, 12, 8, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 1, 12, 3, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 11, 2, 13, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 11, 12, 13, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 1, 2, 3, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 11, 2, 13, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 1, 12, 3, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 11, 2, 13, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 1, 12, 3, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 1, 7, 13, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 11, 7, 3, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 11, 7, 13, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 1, 7, 3, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 6, 2, 3, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 6, 12, 13, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 1, 12, 13, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 11, 2, 3, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 1, 12, 8, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 11, 2, 8, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 6, 2, 8, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 6, 12, 8, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 11, 12, 3, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 1, 7, 8, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 11, 7, 8, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 1, 2, 3, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 11, 12, 13, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 11, 12, 13, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 1, 2, 3, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 1, 7, 3, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 11, 7, 13, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 6, 12, 13, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 6, 2, 3, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 6, 7, 8, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 6, 7, 8, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 1, 2, 3, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 1, 2, 8, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 1, 2, 8, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 1, 7, 3, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 1, 7, 8, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 1, 7, 8, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 1, 7, 13, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 6, 2, 3, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 6, 2, 8, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 1, 7, 3, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 11, 7, 13, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 6, 2, 13, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 6, 12, 3, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 6, 7, 8, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 6, 7, 8, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 6, 7, 13, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 1, 2, 3, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 1, 2, 8, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 1, 2, 8, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 6, 7, 3, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 1, 7, 8, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 1, 7, 13, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 6, 2, 3, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 6, 2, 3, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 6, 2, 8, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 1, 12, 8, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 6, 7, 3, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 6, 12, 13, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 11, 12, 13, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 6, 7, 8, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 6, 2, 13, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 11, 12, 13, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 11, 7, 3, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 11, 7, 3, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 11, 7, 8, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 6, 12, 8, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 1, 2, 13, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 11, 12, 3, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 1, 12, 13, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 11, 2, 3, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 11, 2, 8, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 1, 7, 8, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 11, 7, 8, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 6, 12, 8, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 6, 12, 13, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{0, 11, 7, 8, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{10, 1, 7, 8, 9})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 6, 7, 13, 4})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 6, 7, 3, 14})
	idx++
	PaylineTarzanMapping.Set(idx, []int{5, 6, 7, 8, 4})
}

func NewTarzanMatrix() SlotMatrix {
	m := SlotMatrix{
		List: make([]pb.SiXiangSymbol, ColsTarzanMatrix*RowsTarzanMatrix),
		Cols: ColsTarzanMatrix,
		Rows: RowsTarzanMatrix,
		Size: RowsTarzanMatrix * ColsTarzanMatrix,
	}
	return m
}

func NewTarzanJungleTreasureMatrix() SlotMatrix {
	m := SlotMatrix{
		Cols: ColsTarzanMatrix,
		Rows: RowsTarzanMatrix,
	}
	m.Size = m.Cols * m.Rows
	m.List = make([]pb.SiXiangSymbol, 0, m.Size)
	for sym, val := range TarzanJungleTreasureSymbol {
		for i := 0; i < val.NumOccur; i++ {
			m.List = append(m.List, sym)
		}
	}
	m.TrackFlip = make(map[int]bool)
	return m
}
