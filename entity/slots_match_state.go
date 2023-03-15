package entity

import (
	"time"

	"github.com/ciaolink-game-platform/cgp-common/lib"
	orderedmap "github.com/wk8/go-ordered-map/v2"

	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

const (
	MinPresences = 1
	MaxPresences = 1
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
		Cols:      RowsMatrixGoldPick,
		Rows:      ColsMatrixGoldPick,
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
	for _, symbol := range ListSymbolSiXiangBonusGame {
		sm.List = append(sm.List, symbol)
	}
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

func (sm *SlotMatrix) ForEeachLine(fn func(line int, symbols []pb.SiXiangSymbol)) {
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
		if sm.TrackFlip[id] == false {
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

type UserDataMatch struct {
	RRSpecialGame bool `json:"rr_spec_game"`
}

type SlotsMatchState struct {
	lib.MatchState
	allowSpin     bool // allow user submit new bet
	balanceResult *pb.BalanceResult
	// UserDataMatch  UserDataMatch
	matrix         SlotMatrix
	spreadMatrix   SlotMatrix
	paylines       []*pb.Payline
	WaitSpinMatrix bool

	CurrentSiXiangGame pb.SiXiangGame
	// next game in loop
	NextSiXiangGame pb.SiXiangGame
	bet             *pb.InfoBet
	WinJp           pb.WinJackpot

	MatrixSpecial SlotMatrix
	// ChipsWinInSpecialGame int64
	SpinSymbols      []*pb.SpinSymbol
	EyeSiXiangRemain []pb.SiXiangSymbol
	GemSpin          int // gem using for spin in dragon perl
	// lần quay chắc chắn ra ngọc
	TurnSureSpinEye int
	// Danh sach ngoc tứ linh spin được theo chip bet.
	EyeSiXiangSpined map[int][]pb.SiXiangSymbol
}

func NewSlotsMathState(label *lib.MatchLabel) *SlotsMatchState {
	m := SlotsMatchState{
		MatchState:     lib.NewMathState(label, NewMyPrecense),
		balanceResult:  nil,
		WaitSpinMatrix: false,
		bet: &pb.InfoBet{
			Chips: 0,
		},
		CurrentSiXiangGame: pb.SiXiangGame_SI_XIANG_GAME_NORMAL,
		NextSiXiangGame:    pb.SiXiangGame_SI_XIANG_GAME_NORMAL,
	}

	return &m
}

func (s *SlotsMatchState) GetMatrix() SlotMatrix {
	return s.matrix
}

func (s *SlotsMatchState) SetMatrix(matrix SlotMatrix) {
	s.matrix = matrix
}

func (s *SlotsMatchState) GetSpreadMatrix() SlotMatrix {
	return s.spreadMatrix
}

func (s *SlotsMatchState) SetSpreadMMatrix(matrix SlotMatrix) {
	s.spreadMatrix = matrix
}

func (s *SlotsMatchState) GetPaylines() []*pb.Payline {
	return s.paylines
}

func (s *SlotsMatchState) SetPaylines(paylines []*pb.Payline) {
	s.paylines = paylines
}

func (s *SlotsMatchState) GetBetInfo() *pb.InfoBet {
	return s.bet
}

func (s *SlotsMatchState) SetBetInfo(bet *pb.InfoBet) {
	s.bet = bet
}

func (s *SlotsMatchState) InitNewRound() {
	s.WaitSpinMatrix = false
	s.WinJp = pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED
	s.paylines = nil
	s.spreadMatrix = SlotMatrix{}
}
func (s *SlotsMatchState) IsAllowSpin() bool {
	return s.allowSpin
}

func (s *SlotsMatchState) SetAllowSpin(val bool) {
	s.allowSpin = val
}

func (s *SlotsMatchState) SetUpCountDown(duration time.Duration) {
	s.CountDownReachTime = time.Now().Add(duration)
	s.LastCountDown = -1
}

func (s *SlotsMatchState) ResetBalanceResult() {
	s.SetBalanceResult(nil)
}

func (s *SlotsMatchState) GetBalanceResult() *pb.BalanceResult {
	return s.balanceResult
}

func (s *SlotsMatchState) SetBalanceResult(u *pb.BalanceResult) {
	s.balanceResult = u
}

// func (s *SlotsMatchState) ResetTrackingPlayBonusGame() {
// 	s.trackingPlaySiXiangGame = make(map[pb.SiXiangGame]int)
// }
