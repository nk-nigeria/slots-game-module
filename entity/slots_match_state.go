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

type SixiangMatchState struct {
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

func NewSlotsMathState(label *lib.MatchLabel) *SixiangMatchState {
	m := SixiangMatchState{
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

func (s *SixiangMatchState) GetMatrix() SlotMatrix {
	return s.matrix
}

func (s *SixiangMatchState) SetMatrix(matrix SlotMatrix) {
	s.matrix = matrix
}

func (s *SixiangMatchState) GetSpreadMatrix() SlotMatrix {
	return s.spreadMatrix
}

func (s *SixiangMatchState) SetSpreadMMatrix(matrix SlotMatrix) {
	s.spreadMatrix = matrix
}

func (s *SixiangMatchState) GetPaylines() []*pb.Payline {
	return s.paylines
}

func (s *SixiangMatchState) SetPaylines(paylines []*pb.Payline) {
	s.paylines = paylines
}

func (s *SixiangMatchState) GetBetInfo() *pb.InfoBet {
	return s.bet
}

func (s *SixiangMatchState) SetBetInfo(bet *pb.InfoBet) {
	s.bet = bet
}

func (s *SixiangMatchState) InitNewRound() {
	s.WaitSpinMatrix = false
	s.WinJp = pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED
	s.paylines = nil
	s.spreadMatrix = SlotMatrix{}
}
func (s *SixiangMatchState) IsAllowSpin() bool {
	return s.allowSpin
}

func (s *SixiangMatchState) SetAllowSpin(val bool) {
	s.allowSpin = val
}

func (s *SixiangMatchState) SetUpCountDown(duration time.Duration) {
	s.CountDownReachTime = time.Now().Add(duration)
	s.LastCountDown = -1
}

func (s *SixiangMatchState) ResetBalanceResult() {
	s.SetBalanceResult(nil)
}

func (s *SixiangMatchState) GetBalanceResult() *pb.BalanceResult {
	return s.balanceResult
}

func (s *SixiangMatchState) SetBalanceResult(u *pb.BalanceResult) {
	s.balanceResult = u
}

// func (s *SlotsMatchState) ResetTrackingPlayBonusGame() {
// 	s.trackingPlaySiXiangGame = make(map[pb.SiXiangGame]int)
// }
