package entity

import pb "github.com/ciaolink-game-platform/cgp-common/proto"

const (
	Row_1 = 0
	Row_2 = 1
	Row_3 = 2
	Row_4 = 3
	Row_5 = 4

	Col_1 = 0
	Col_2 = 1
	Col_3 = 2
	Col_4 = 3
	Col_5 = 4
)

type SlotMatrix struct {
	List      []pb.SiXiangSymbol
	Cols      int
	Rows      int
	Size      int
	TrackFlip map[int]bool
}

func NewSiXiangMatrixNormal() SlotMatrix {
	sm := SlotMatrix{
		List:      make([]pb.SiXiangSymbol, RowsMatrix*ColsMatrix, RowsMatrix*ColsMatrix),
		Cols:      ColsMatrix,
		Rows:      RowsMatrix,
		Size:      0,
		TrackFlip: map[int]bool{},
	}
	sm.Size = sm.Cols * sm.Rows
	return sm
}

func NewMatrixBonusGame() SlotMatrix {
	sm := SlotMatrix{
		List:      make([]pb.SiXiangSymbol, 0, len(ListSymbolBonusGame)),
		Cols:      len(ListSymbolBonusGame) / 2,
		Rows:      2,
		TrackFlip: map[int]bool{},
	}
	sm.Size = sm.Cols * sm.Rows
	for symbol := range ListSymbolBonusGame {
		sm.List = append(sm.List, symbol)
	}
	return sm
}

func NewMatrixLuckyDraw() SlotMatrix {
	sm := SlotMatrix{
		List:      make([]pb.SiXiangSymbol, 0, RowsMatrix*ColsMatrix),
		Cols:      ColsMatrix,
		Rows:      RowsMatrix,
		TrackFlip: map[int]bool{},
	}
	for symbol, info := range ListSymbolLuckyDraw {
		for i := 0; i < info.NumOccur; i++ {
			sm.List = append(sm.List, symbol)
		}
	}
	sm.Size = sm.Cols * sm.Rows
	return sm
}

func NewMatrixDragonPearl() SlotMatrix {
	sm := SlotMatrix{
		List:      make([]pb.SiXiangSymbol, 0, RowsMatrix*ColsMatrix),
		Cols:      ColsMatrix,
		Rows:      RowsMatrix,
		TrackFlip: map[int]bool{},
	}
	for symbol, v := range ListSymbolDragonPearl {
		for i := 0; i < v.NumOccur; i++ {
			sm.List = append(sm.List, symbol)
		}
	}
	sm.Size = sm.Cols * sm.Rows
	return sm
}

func NewMatrixGoldPick() SlotMatrix {
	sm := SlotMatrix{
		List:      make([]pb.SiXiangSymbol, 0, RowsMatrixGoldPick*ColsMatrixGoldPick),
		Cols:      RowsMatrixGoldPick,
		Rows:      ColsMatrixGoldPick,
		TrackFlip: map[int]bool{},
	}
	for symbol, v := range ListSymbolGoldPick {
		for i := 0; i < v.NumOccur; i++ {
			sm.List = append(sm.List, symbol)
		}
	}
	sm.Size = sm.Cols * sm.Rows
	return sm
}

// x4 END
// x3 x4 END
// x2 x3 x4 END
// x2 x3 X4 END
// x2 x2 x3 x3 x4
func NewMatrixRapidPay() SlotMatrix {
	sm := SlotMatrix{
		List:      make([]pb.SiXiangSymbol, 0, RowsMatrixRapidPay*ColsMatrixRapidPay),
		Cols:      ColsMatrixRapidPay,
		Rows:      RowsMatrixRapidPay,
		TrackFlip: map[int]bool{},
	}
	// x4 END
	sm.List = append(sm.List,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X4,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_END,
	)
	sm.List = append(sm.List, SliceRepeat(3, pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED)...)
	// x3 x4 END
	sm.List = append(sm.List,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X3,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X4,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_END,
	)
	sm.List = append(sm.List, SliceRepeat(2, pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED)...)
	// x2 x3 x4 END
	sm.List = append(sm.List,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X2,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X3,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X4,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_END,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED,
	)
	// x2 x3 X4 END
	sm.List = append(sm.List,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X2,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X3,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X4,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_END,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED,
	)
	// x2 x2 x3 x3 x4
	sm.List = append(sm.List,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X2,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X2,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X3,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X3,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X4,
	)
	sm.Size = sm.Cols * sm.Rows
	return sm
}

func NewMatrixSiXiangBonus() SlotMatrix {
	sm := SlotMatrix{
		List:      make([]pb.SiXiangSymbol, 0, len(ListSymbolSiXiangBonusGame)),
		Cols:      len(ListSymbolBonusGame) / 2,
		Rows:      2,
		TrackFlip: map[int]bool{},
	}
	sm.Size = sm.Cols * sm.Rows
	sm.List = append(sm.List, ListSymbolSiXiangBonusGame...)
	return sm
}

func (sm *SlotMatrix) RowCol(id int) (int, int) {
	if sm.Cols == 0 || sm.Rows == 0 {
		return 0, 0
	}
	row := id / sm.Cols
	col := id - row*sm.Cols
	return row, col
}

func (sm *SlotMatrix) Reset() {
	sm.List = make([]pb.SiXiangSymbol, 0)
	sm.Cols = 0
	sm.Rows = 0
	sm.TrackFlip = make(map[int]bool)
}

func (sm *SlotMatrix) ForEeach(fn func(idx, row, col int, symbol pb.SiXiangSymbol)) {
	row := 0
	col := 0
	cols := sm.Cols
	for idx, symbol := range sm.List {
		if idx != 0 {
			if col%cols == 0 {
				row++
				col = 0
			}
		}
		fn(idx, row, col, symbol)
		col++
	}
}

func (sm *SlotMatrix) ForEachLine(fn func(line int, symbols []pb.SiXiangSymbol)) {
	list := make([]pb.SiXiangSymbol, 0)
	sm.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if idx != 0 && col == 0 {
			fn(row-1, list)
			list = make([]pb.SiXiangSymbol, 0)
		}
		list = append(list, symbol)
	})
	fn(sm.Rows-1, list)
}

func (sm *SlotMatrix) ForEachCol(fn func(col int, symbols []pb.SiXiangSymbol)) {
	for col := 0; col < sm.Cols; col++ {
		ml := make([]pb.SiXiangSymbol, 0, sm.Rows)
		for row := 0; row < sm.Rows; row++ {
			idx := row*sm.Cols + col
			ml = append(ml, sm.List[idx])
		}
		fn(col, ml)
	}
}

func (sm *SlotMatrix) ListFromIndexs(indexs []int) []pb.SiXiangSymbol {
	size := len(sm.List)
	list := make([]pb.SiXiangSymbol, 0, len(indexs))
	for _, idx := range indexs {
		if idx >= size {
			continue
		}
		list = append(list, sm.List[idx])
	}
	return list
}

func (sm *SlotMatrix) ToPbSlotMatrix() *pb.SlotMatrix {
	// matrix, cols,row
	pbSl := &pb.SlotMatrix{
		Rows: int32(sm.Rows),
		Cols: int32(sm.Cols),
	}
	pbSl.Lists = make([]pb.SiXiangSymbol, pbSl.Rows*pbSl.Cols)
	copy(pbSl.Lists, sm.List) // deep copy
	return pbSl
}

func (sm *SlotMatrix) RandomSymbolNotFlip(randomFn func(min, max int) int) (int, pb.SiXiangSymbol) {
	listIdNotFlip := make([]int, 0)
	for id := range sm.List {
		if !sm.TrackFlip[id] {
			listIdNotFlip = append(listIdNotFlip, id)
		}
	}
	id := randomFn(0, len(listIdNotFlip))
	idInList := listIdNotFlip[id]
	symbol := sm.List[idInList]
	return idInList, symbol
}

func (sm *SlotMatrix) Flip(idx int) pb.SiXiangSymbol {
	sm.TrackFlip[idx] = true
	return sm.List[idx]
}

func (sm *SlotMatrix) IsPayline(paylineIndex []int) ([]int, bool) {
	if len(paylineIndex) == 0 || len(sm.List) == 0 {
		return nil, false
	}
	firstSymbol := sm.List[paylineIndex[0]]
	validPaylineIndex := make([]int, 0)
	for _, idx := range paylineIndex {
		sym := sm.List[idx]
		if firstSymbol == sym || sym == pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD || sym == pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN {
			validPaylineIndex = append(validPaylineIndex, idx)
			continue
		}
		break

	}
	if len(validPaylineIndex) >= sm.Cols {
		return validPaylineIndex, true
	}
	return nil, false
}
