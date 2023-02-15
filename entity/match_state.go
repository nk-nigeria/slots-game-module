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

type SlotsMatchState struct {
	lib.MatchState
	allowBet      bool // allow user submit new bet
	balanceResult *pb.BalanceResult

	userBets    map[string][]*pb.InfoBet
	colorBets   map[int]*pb.InfoBet
	historyRoll [][]int32 // list id color
	// deck      *pb.Deck
}

func NewSlotsMathState(label *lib.MatchLabel) SlotsMatchState {
	m := SlotsMatchState{
		MatchState:    lib.NewMathState(label, NewMyPrecense),
		balanceResult: nil,
		userBets:      make(map[string][]*pb.InfoBet),
		colorBets:     make(map[int]*pb.InfoBet),
		// deck:                nil,
	}
	return m
}

func (s *SlotsMatchState) InitNewMatch() {
	for k := range s.userBets {
		delete(s.userBets, k)
	}
	for k := range s.colorBets {
		delete(s.colorBets, k)
	}
	s.balanceResult = nil

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
