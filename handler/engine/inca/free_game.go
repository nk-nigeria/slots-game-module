package inca

import (
	"github.com/nakamaFramework/cgb-slots-game-module/entity"
	"github.com/nakamaFramework/cgp-common/lib"
	pb "github.com/nakamaFramework/cgp-common/proto"
)

var _ lib.Engine = &freeGame{}

type freeGame struct {
	normal
}

func NewFreeGame(randomIntFn func(int, int) int) lib.Engine {
	e := &freeGame{}
	if randomIntFn != nil {
		e.randomFn = randomIntFn
	} else {
		e.randomFn = entity.RandomInt
	}
	return e
}

// Finish implements lib.Engine.
func (e *freeGame) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	paylines := s.Paylines()
	totalWin := int64(0)
	for _, payline := range paylines {
		payline.Chips = s.Bet().Chips * int64(payline.Rate) / 20
		totalWin += payline.Chips
	}
	{
		numScatter := e.countScatterByCol(s.Matrix)
		ratio := entity.IncaRatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER][int32(numScatter)]
		totalWin += int64(ratio) * s.Bet().Chips / 20
	}
	s.ChipStat.AddChipWin(s.CurrentSiXiangGame, totalWin)
	s.NextSiXiangGame = e.GetNextSiXiangGame(s)
	slotDesk := &pb.SlotDesk{
		GameReward: &pb.GameReward{
			ChipsWin:            s.ChipStat.ChipWin(s.CurrentSiXiangGame),
			TotalChipsWinByGame: s.ChipStat.TotalChipWin(s.CurrentSiXiangGame),
		},
		ChipsMcb:           s.Bet().Chips,
		CurrentSixiangGame: s.CurrentSiXiangGame,
		NextSixiangGame:    s.NextSiXiangGame,
		Matrix:             s.Matrix.ToPbSlotMatrix(),
		SpreadMatrix:       s.WildMatrix.ToPbSlotMatrix(),
		Paylines:           paylines,
		IsFinishGame:       s.NumSpinLeft <= 0,
		NumSpinLeft:        int64(s.NumSpinLeft),
		BetLevels:          entity.BetLevels[:],
		GameConfig:         s.GameConfig.GameConfig,
	}
	if slotDesk.IsFinishGame {
		s.GameConfig.GameConfig = &pb.GameConfig{}
	}
	s.LastResult = slotDesk
	return slotDesk, nil
}

// NewGame implements lib.Engine.
func (e *freeGame) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.NumSpinLeft <= 0 {
		s.NumSpinLeft = 15
	}
	return matchState, nil
}

func (e *freeGame) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	var err error
	for {
		_, err = e.normal.Process(s)
		if err != nil {
			return s, err
		}
		numScatterSeq := e.countScatterByCol(s.Matrix)
		if numScatterSeq >= 3 {
			// s.NumSpinLeft += 15
			// not allow 3 scatter in freespin
			s.NumSpinLeft += 1
			continue
		}
		break
	}
	return s, err
}

func (e *freeGame) GetNextSiXiangGame(s *entity.SlotsMatchState) pb.SiXiangGame {
	if s.NumSpinLeft <= 0 {
		return pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	}
	return s.CurrentSiXiangGame
}
