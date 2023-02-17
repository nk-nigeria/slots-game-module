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

const (
	MaxColMatrix = 5
	MaxRowMatix  = 3
)

type SlotMatrix [MaxRowMatix][MaxColMatrix]int

var ListSymbol = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 100, 200, 300}

type SlotsMatchState struct {
	lib.MatchState
	allowBet      bool // allow user submit new bet
	balanceResult *pb.BalanceResult

	matrix         SlotMatrix
	WaitSpinMatrix bool
}

func NewSlotsMathState(label *lib.MatchLabel) *SlotsMatchState {
	m := SlotsMatchState{
		MatchState:     lib.NewMathState(label, NewMyPrecense),
		balanceResult:  nil,
		WaitSpinMatrix: false,
	}
	return &m
}

func (s *SlotsMatchState) GetMatrix() (SlotMatrix, int, int) {
	return s.matrix, MaxColMatrix, MaxRowMatix
}

func (s *SlotsMatchState) SetMatrix(matrix SlotMatrix) {
	s.matrix = matrix
}

func (s *SlotsMatchState) InitNewMatch() {

}

func (s *SlotsMatchState) IsAllowBet() bool {
	return s.allowBet
}

func (s *SlotsMatchState) SetAllowBet(val bool) {
	s.allowBet = val
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
