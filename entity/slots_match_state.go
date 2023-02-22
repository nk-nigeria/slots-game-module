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
	MaxColMatrix = 5
	MaxRowMatix  = 3
)

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
	List []pb.SiXiangSymbol
	Cols int
	Rows int
}

func NewSlotMatrix() SlotMatrix {
	sm := SlotMatrix{
		List: make([]pb.SiXiangSymbol, MaxRowMatix*MaxColMatrix, MaxRowMatix*MaxColMatrix),
		Cols: MaxColMatrix,
		Rows: MaxRowMatix,
	}

	return sm
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
		Rows:  int32(sm.Rows),
		Cols:  int32(sm.Cols),
		Lists: sm.List[:], // using trick [:] for deep copy list
	}
	return pbSl
}

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

type SlotsMatchState struct {
	lib.MatchState
	allowSpin     bool // allow user submit new bet
	balanceResult *pb.BalanceResult

	matrix                  SlotMatrix
	spreadMatrix            SlotMatrix
	paylines                []*pb.Payline
	WaitSpinMatrix          bool
	AutoSpin                bool
	playSiXiangGameTracking map[pb.SiXiangGame]int
	// current game
	SiXiangGame pb.SiXiangGame
	// next game in loop
	NextSiXiangGame pb.SiXiangGame
	bet             *pb.InfoBet
}

func NewSlotsMathState(label *lib.MatchLabel) *SlotsMatchState {
	m := SlotsMatchState{
		MatchState:     lib.NewMathState(label, NewMyPrecense),
		balanceResult:  nil,
		WaitSpinMatrix: false,
	}
	m.matrix = NewSlotMatrix()
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

func (s *SlotsMatchState) InitNewMatch() {
	s.AutoSpin = true
	s.WaitSpinMatrix = false

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

func (s *SlotsMatchState) IsWinSiXangBonusGame() bool {
	return s.playSiXiangGameTracking[pb.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL]&
		s.playSiXiangGameTracking[pb.SiXiangGame_SI_XIANG_GAME_LUCKDRAW]&
		s.playSiXiangGameTracking[pb.SiXiangGame_SI_XIANG_GAME_GOLDPICK]&
		s.playSiXiangGameTracking[pb.SiXiangGame_SI_XIANG_GAME_RAPIDPAY] != 0
}

func (s *SlotsMatchState) AddTrackingPlayBonusGame(siXiangGame pb.SiXiangGame) {
	num := s.playSiXiangGameTracking[siXiangGame]
	num++
	s.playSiXiangGameTracking[siXiangGame] = num
}

func (s *SlotsMatchState) ResetdTrackingPlayBonusGame() {
	s.playSiXiangGameTracking = make(map[pb.SiXiangGame]int)
}
