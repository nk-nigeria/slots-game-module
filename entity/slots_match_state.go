package entity

import (
	"time"

	"github.com/ciaolink-game-platform/cgp-common/lib"

	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

const (
	MinPresences = 1
	MaxPresences = 1
)

type UserDataMatch struct {
	RRSpecialGame bool `json:"rr_spec_game"`
}

type SlotsMatchState struct {
	lib.MatchState
	allowSpin     bool // allow user submit new bet
	balanceResult *pb.BalanceResult
	// UserDataMatch  UserDataMatch
	Matrix         SlotMatrix
	WildMatrix     SlotMatrix
	paylines       []*pb.Payline
	WaitSpinMatrix bool

	CurrentSiXiangGame pb.SiXiangGame
	// next game in loop
	NextSiXiangGame pb.SiXiangGame
	bet             *pb.InfoBet
	WinJp           pb.WinJackpot

	MatrixSpecial SlotMatrix
	// ChipsWinInSpecialGame int64
	SpinSymbols []*pb.SpinSymbol
	GemSpin     int // gem using for spin in dragon perl
	// lần quay chắc chắn ra ngọc
	TurnSureSpin           int
	CollectionSymbolRemain []pb.SiXiangSymbol
	// Danh sach ngoc tứ linh spin được theo chip bet.
	CollectionSymbol map[int]map[pb.SiXiangSymbol]int

	// tarzan
	// List idx of free symbol index
	TrackIndexFreeSpinSymbol     map[int]bool
	ChipWinByGame                map[pb.SiXiangGame]int64
	LineWinByGame                map[pb.SiXiangGame]int
	CountLineCrossFreeSpinSymbol int
	// ngoc rung xanh
	PerlGreenForest int
	// chip tich luy
	ChipsBonus int64
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
		CollectionSymbol:   make(map[int]map[pb.SiXiangSymbol]int, 0),
		ChipWinByGame:      make(map[pb.SiXiangGame]int64, 0),
		LineWinByGame:      make(map[pb.SiXiangGame]int, 0),
	}

	return &m
}

func (s *SlotsMatchState) SetMatrix(matrix SlotMatrix) {
	s.Matrix = matrix
}

func (s *SlotsMatchState) SetWildMatrix(matrix SlotMatrix) {
	s.WildMatrix = matrix
}

func (s *SlotsMatchState) Paylines() []*pb.Payline {
	return s.paylines
}

func (s *SlotsMatchState) SetPaylines(paylines []*pb.Payline) {
	s.paylines = paylines
}

func (s *SlotsMatchState) Bet() *pb.InfoBet {
	return s.bet
}

func (s *SlotsMatchState) SetBetInfo(bet *pb.InfoBet) {
	s.bet = bet
}

func (s *SlotsMatchState) InitNewRound() {
	s.WaitSpinMatrix = false
	s.WinJp = pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED
	s.paylines = nil
	s.WildMatrix = SlotMatrix{}
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

//	func (s *SlotsMatchState) ResetTrackingPlayBonusGame() {
//		s.trackingPlaySiXiangGame = make(map[pb.SiXiangGame]int)
//	}
func (s *SlotsMatchState) AddCollectionSymbol(chipMcb int, sym pb.SiXiangSymbol) {
	collection, exist := s.CollectionSymbol[chipMcb]
	if !exist {
		s.CollectionSymbol[chipMcb] = make(map[pb.SiXiangSymbol]int, 0)
		collection = s.CollectionSymbol[chipMcb]
	}
	num := collection[sym]
	num++
	collection[sym] = num
	s.CollectionSymbol[chipMcb] = collection
}

func (s *SlotsMatchState) CollectionSymbolToSlice(chipMcb int) []pb.SiXiangSymbol {
	collection, exist := s.CollectionSymbol[chipMcb]
	ml := make([]pb.SiXiangSymbol, 0, len(s.CollectionSymbol))
	if !exist {
		return ml
	}
	for k := range collection {
		ml = append(ml, k)
	}
	return ml
}

func (s *SlotsMatchState) SizeCollectionSymbol(chipMcb int) int {
	return len(s.CollectionSymbol[chipMcb])
}
